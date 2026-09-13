<?php
/**
 * Roundcube autologin endpoint for OpenCloud Webmail extension.
 *
 * Receives HMAC-SHA256 signed tokens from the browser extension and
 * authenticates the user into Roundcube via its internal login API.
 *
 * Query parameters:
 *   data - base64url-encoded JSON payload {sub, accountId, email, imapPass, exp}
 *   sig  - base64url(HMAC-SHA256(data + "." + ts, SHARED_SECRET))
 *   ts   - Unix timestamp (seconds) when the token was created
 *
 * Install: place this file in the Roundcube web root alongside index.php.
 * Configure SHARED_SECRET below (must match the OpenCloud apps.yaml value).
 */

define('SHARED_SECRET', getenv('ROUNDCUBE_AUTOLOGIN_SECRET'));
define('MAX_TOKEN_AGE', 120); // seconds

function base64url_decode_safe(string $input): string|false {
    $padded = str_pad(strtr($input, '-_', '+/'), (int)(ceil(strlen($input) / 4) * 4), '=');
    return base64_decode($padded, true);
}

function fail(string $message, int $code = 403): never {
    http_response_code($code);
    header('Content-Type: text/plain');
    echo $message;
    exit;
}

// --- Validate request ---

$data = $_GET['data'] ?? '';
$sig  = $_GET['sig']  ?? '';
$ts   = $_GET['ts']   ?? '';

if (!$data || !$sig || !$ts) {
    fail('Missing parameters', 400);
}

// Check timestamp freshness
$tsInt = (int)$ts;
if (abs(time() - $tsInt) > MAX_TOKEN_AGE) {
    fail('Token expired');
}

// Verify HMAC signature
$message = $data . '.' . $ts;
$expectedSig = rtrim(strtr(base64_encode(
    hash_hmac('sha256', $message, SHARED_SECRET, true)
), '+/', '-_'), '=');

if (!hash_equals($expectedSig, $sig)) {
    fail('Invalid signature');
}

// Decode payload
$json = base64url_decode_safe($data);
if ($json === false) {
    fail('Invalid data encoding', 400);
}

$payload = json_decode($json, true);
if (!$payload || !isset($payload['email'], $payload['imapPass'], $payload['exp'])) {
    fail('Invalid payload', 400);
}

// Check expiry in payload
if ($payload['exp'] < time()) {
    fail('Token expired (payload)');
}

// --- Roundcube login ---

define('INSTALL_PATH', realpath(__DIR__ . '/..') . '/');

require_once INSTALL_PATH . 'program/include/iniset.php';

$rcmail = rcmail::get_instance(0, 'larry');

$email = strtolower(trim($payload['email']));

$at = strrpos($email, '@');
if ($at === false) {
    fail('Invalid email address', 400);
}

$domain = substr($email, $at + 1);

$imapHosts = $rcmail->config->get('opencloud_imap_domains', []);

if (!isset($imapHosts[$domain])) {
    fail('Unsupported email domain');
}

$auth = $rcmail->login(
    $payload['email'],
    $payload['imapPass'],
    $imapHosts[$domain],
    false
);

if (!$auth) {
    fail('IMAP login failed');
}

// Set session cookie and redirect
$rcmail->session->set_auth_cookie();

header('Location: ./');
exit;
