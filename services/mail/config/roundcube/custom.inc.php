<?php

$config['x_frame_options'] = false;
$config['session_lifetime'] = 600;

$config['imap_host'] = [
    'ssl://' . getenv('HOMELAB_MAIL_HOST') . ':993' => getenv('HOMELAB_MAIL_NAME') ?: 'Homelab Mail',
    'ssl://' . getenv('WORK_MAIL_HOST') . ':993' => getenv('WORK_MAIL_NAME') ?: 'Work Mail',
    'ssl://imap.gmail.com:993' => 'Gmail',
];

$config['smtp_host'] = [
    getenv('HOMELAB_MAIL_HOST') => 'ssl://' . getenv('HOMELAB_MAIL_HOST') . ':465',
    getenv('WORK_MAIL_HOST') => 'ssl://' . getenv('WORK_MAIL_HOST') . ':465',
    'imap.gmail.com' => 'ssl://smtp.gmail.com:465',
];

$config['smtp_user'] = '%u';
$config['smtp_pass'] = '%p';
