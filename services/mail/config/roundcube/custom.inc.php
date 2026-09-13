<?php

$config['x_frame_options'] = false;
$config['session_lifetime'] = 600;


// -----------------------------------------------------------------------------
// Mail server configuration
// -----------------------------------------------------------------------------

// Homelab
$homelabHost   = getenv('HOMELAB_MAIL_HOST');
$homelabName   = getenv('HOMELAB_MAIL_NAME');
$homelabDomain = getenv('HOMELAB_MAIL_DOMAIN');

// Work
$workHost   = getenv('WORK_MAIL_HOST');
$workName   = getenv('WORK_MAIL_NAME');
$workDomain = getenv('WORK_MAIL_DOMAIN');


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

$config['smtp_host'] = [
    $homelabHost     => 'ssl://' . $homelabHost . ':465',
    'imap.gmail.com' => 'ssl://smtp.gmail.com:465',
    $workHost        => 'ssl://' . $workHost . ':465',
];

$config['smtp_user'] = '%u';
$config['smtp_pass'] = '%p';


// -----------------------------------------------------------------------------
// OpenCloud autologin domain -> IMAP routing
// -----------------------------------------------------------------------------

$config['opencloud_imap_domains'] = [
    strtolower($homelabDomain) => 'ssl://' . $homelabHost . ':993',
    'gmail.com'                 => 'ssl://imap.gmail.com:993',
    strtolower($workDomain)    => 'ssl://' . $workHost . ':993',
];
