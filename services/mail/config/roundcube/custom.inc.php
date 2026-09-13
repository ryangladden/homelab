<?php

$config['x_frame_options'] = false;
$config['session_lifetime'] = 600;

// -----------------------------------------------------------------------------
// Mail server configuration
// -----------------------------------------------------------------------------

// Homelab
$homelabHost   = trim(getenv('HOMELAB_MAIL_HOST') ?: '');
$homelabName   = trim(getenv('HOMELAB_MAIL_NAME') ?: 'Homelab');
$homelabDomain = strtolower(trim(getenv('HOMELAB_MAIL_DOMAIN') ?: ''));

// Work
$workHost   = trim(getenv('WORK_MAIL_HOST') ?: '');
$workName   = trim(getenv('WORK_MAIL_NAME') ?: 'Work');
$workDomain = strtolower(trim(getenv('WORK_MAIL_DOMAIN') ?: ''));

// Required configuration
if (!$homelabHost || !$homelabDomain) {
    throw new RuntimeException(
        'HOMELAB_MAIL_HOST and HOMELAB_MAIL_DOMAIN must be configured'
    );
}

if (!$workHost || !$workDomain) {
    throw new RuntimeException(
        'WORK_MAIL_HOST and WORK_MAIL_DOMAIN must be configured'
    );
}


// -----------------------------------------------------------------------------
// IMAP
// -----------------------------------------------------------------------------

$config['imap_host'] = [
    'ssl://' . $homelabHost . ':993' => $homelabName,
    'ssl://imap.gmail.com:993'       => 'Google',
    'ssl://' . $workHost . ':993'    => $workName,
];


// -----------------------------------------------------------------------------
// SMTP
// -----------------------------------------------------------------------------
//
// Roundcube 1.7.4 supports per-IMAP-host SMTP configuration.
// Keys MUST be the normalized IMAP hostname, without ssl:// or port.
//

$config['smtp_host'] = [
    $homelabHost     => 'ssl://' . $homelabHost . ':465',
    'imap.gmail.com' => 'ssl://smtp.gmail.com:465',
    $workHost        => 'ssl://' . $workHost . ':465',
];

// Reuse the IMAP username/password for SMTP authentication.
$config['smtp_user'] = [
    $homelabHost     => '%u',
    'imap.gmail.com' => '%u',
    $workHost        => '%u',
];

$config['smtp_pass'] = [
    $homelabHost     => '%p',
    'imap.gmail.com' => '%p',
    $workHost        => '%p',
];


// -----------------------------------------------------------------------------
// OpenCloud autologin domain -> IMAP routing
// -----------------------------------------------------------------------------

$config['opencloud_imap_domains'] = [
    $homelabDomain => 'ssl://' . $homelabHost . ':993',
    'gmail.com'    => 'ssl://imap.gmail.com:993',
    $workDomain    => 'ssl://' . $workHost . ':993',
];
