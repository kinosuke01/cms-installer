<?php
$pwd = __DIR__;
require_once("$pwd/init.php");

function errf($line, $desc, $case, $key, $want, $got) {
    return sprintf(
        "%s:%s %s:%s %s wrong. want=%s got=%s",
        __FILE__, $line, $desc, $case, $key, $want, $got
    );
}

function setup($pwd = __DIR__) {
    $cmds = array(
        "rm -Rf /var/www/html/*",
        "if [ ! -d $pwd/tmp ]; then mkdir $pwd/tmp; fi",
        "rm -Rf $pwd/tmp/*",
        "cp -Rf $pwd/testdata/* /var/www/html/.",
    );

    foreach ($cmds as $cmd) {
        exec($cmd);
    }
}

$funcs = array(
    'CmsArchive#download' => function($desc, $pwd = __DIR__) {
        $errs = array();
        $tt = array(
            array(
                'case' => 'file_exists',
                'url' => 'http://localhost/cms.zip',
                'expected_error' => '',
            ),
            array(
                'case' => 'file_not_exists',
                'url' => 'http://localhost/404.zip',
                'expected_error' => 'DOWNLOAD_ERROR STATUS_CODE=404',
            ),
            array(
                'case' => 'not_connected',
                'url' => 'https://localhost/cms.zip',
                'expected_error' => 'DOWNLOAD_ERROR',
            ),
        );
        foreach ($tt as $tc) {
            $cmsArchive = new CmsArchive();
            $cmsArchive->url = $tc['url'];
            $cmsArchive->basePath = "$pwd/tmp/";

            $errMsg = '';
            try {
                $cmsArchive->download();
            } catch(Exception $e) {
                $errMsg = $e->getMessage();
            }

            if ($tc['expected_error'] != $errMsg) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'error', $tc['expected_error'], $errMsg);
            }
        }
        return $errs;
    },
    'CmsArchive#extract' => function($desc, $pwd = __DIR__) {
        $errs = array();
        $tt = array(
            array(
                'case' => 'file_not_exists',
                'before' => '',
                'download_path' => '404.zip',
                'expected_error' => 'ARCHIVE_FILE_OPEN_ERROR',
                'expected_extraction_result' => false,
            ),
            array(
                'case' => 'not_zip_file_exists',
                'before' => "echo dummyText > $pwd/tmp/not.zip",
                'download_path' => 'not.zip',
                'expected_error' => 'ARCHIVE_FILE_OPEN_ERROR',
                'expected_extraction_result' => false,
            ),
            array(
                'case' => 'zip_file_exists',
                'before' => "cp $pwd/testdata/cms.zip $pwd/tmp/cms.zip",
                'download_path' => 'cms.zip',
                'expected_error' => '',
                'expected_extraction_result' => true,
            )
        );
        foreach($tt as $tc) {
            if ($tc['before'] !== '') {
                exec($tc['before']);
            }
            $cmsArchive = new CmsArchive();
            $cmsArchive->basePath = "$pwd/tmp/";
            $cmsArchive->downloadPath = $tc['download_path'];

            $errMsg = '';
            try {
                $cmsArchive->extract();
            } catch(Exception $e) {
                $errMsg = $e->getMessage();
            }

            if ($tc['expected_error'] != $errMsg) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'error', $tc['expected_error'], $errMsg);
            }

            $extractionResult = is_dir("$pwd/tmp/cms");
            if ($tc['expected_extraction_result'] != $extractionResult) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'error', $tc['extraction_result'], $extractionResult);
            }
        }
        return $errs;
    },
    'CmsArchive#place' => function($desc, $pwd = __DIR__) {
        $errs = array();
        $tt = array(
            array(
                'case' => 'dir_not_exists',
                'before' => "rm -Rf $pwd/tmp/*",
                'extracted_dir' => 'cms',
                'expected_error' => 'PLACE_ERROR(CMS_FILES_PLACE)',
                'expected_file_exists' => false,
                'expected_zip_exists' => false,
            ),
            array(
                'case' => 'dir_exists',
                'before' => "rm -Rf $pwd/tmp/*; cp -R $pwd/testdata/. $pwd/tmp/.",
                'extracted_dir' => 'cms',
                'expected_error' => '',
                'expected_file_exists' => true,
                'expected_zip_exists' => false,
            )
        );
        foreach($tt as $tc) {
            if ($tc['before'] !== '') {
                exec($tc['before']);
            }
            $cmsArchive = new CmsArchive();
            $cmsArchive->basePath = "$pwd/tmp/";
            $cmsArchive->extractedDir = $tc['extracted_dir'];
            $cmsArchive->downloadPath = 'cms.zip';

            $errMsg = '';
            try {
                $cmsArchive->place();
            } catch(Exception $e) {
                $errMsg = $e->getMessage();
            }

            if ($tc['expected_error'] != $errMsg) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'error', $tc['expected_error'], $errMsg);
            }

            $fileExists = is_file("$pwd/tmp/index.html");
            if ($tc['expected_file_exists'] != $fileExists) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'file_exists', $tc['expected_file_exists'], $fileExists);
            }

            $zipExists = is_file("$pwd/tmp/cms.zip");
            if ($tc['expected_zip_exists'] != $zipExists) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'zip_exists', $tc['expected_zip_exists'], $zipExists);
            }
        }
        return $errs;
    },
    'Authenticator#exec' => function($desc) {
        $errs = array();
        $tt = array(
            array(
                'case' => 'invalid_token',
                'itoken' => 'abcde',
                'expired_at' => '1634292000', // '2021-10-15T19:00:00+09:00
                'token' => '12345',
                'now' => '1634288400',
                'expected_error' => 'AUTH_ERROR',
            ),
            array(
                'case' => 'expired',
                'itoken' => 'abcde',
                'expired_at' => '1634292000',
                'token' => 'abcde',
                'now' => '1634292000',
                'expected_error' => 'AUTH_ERROR',
            ),
            array(
                'case' => 'success',
                'itoken' => 'abcde',
                'expired_at' => '1634292000',
                'token' => 'abcde',
                'now' => '1634291999',
                'expected_error' => '',
            ),
        );
        foreach ($tt as $tc) {
            $authenticator = new Authenticator();
            $authenticator->expiredAt = $tc['expired_at'];
            $authenticator->token = $tc['itoken'];

            $errMsg = "";
            try {
                $authenticator->exec($tc['token'], $tc['now']);
            } catch(Exception $e) {
                $errMsg = $e->getMessage();
            }

            if ($tc['expected_error'] !== $errMsg) {
                $errs[] = errf(__LINE__, $desc, $tc['case'], 'error', $tc['expected_error'], $errMsg);
            }
        }
        return $errs;
    },
);

$errs = array();
foreach ($funcs as $desc => $func) {
    setup();
    $errs = array_merge($errs, $func($desc));
}
if (empty($errs)) {
    echo "TEST OK\n";
    exit(0);
} else {
    echo implode("\n", $errs);
    echo "\n";
    exit(1);
}
