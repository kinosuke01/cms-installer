package cmsinit

// TODO sync with init.php
const php string = `<?php
const ARCHIVE_URL   = 'ARCHIVE_URL_PLACEHOLDER';
const EXTRACTED_DIR = 'EXTRACTED_DIR_PLACEHOLDER';
const TOKEN         = 'TOKEN_PLACEHOLDER';
const EXPIRED_AT    = 'EXPIRED_AT_PLACEHOLDER';

function path_join( $base, $path ) {
    return rtrim( $base, '/' ) . '/' . ltrim( $path, '/' );
}

class CmsArchive
{
	public $url = ARCHIVE_URL;
	public $extractedDir = EXTRACTED_DIR;
	public $basePath = './';
	public $downloadPath = 'archive.zip';

	public function download()
	{
		$url = $this->url;
		$filePath = path_join($this->basePath, $this->downloadPath);

		$ch = curl_init($url);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($ch, CURLOPT_BINARYTRANSFER, true);
		curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10);
		curl_setopt($ch, CURLOPT_TIMEOUT, 10);
		$res = curl_exec($ch);
		$statusCode = (int) curl_getinfo($ch, CURLINFO_RESPONSE_CODE);
		curl_close($ch);

		if ($res === false) {
			throw new Exception('DOWNLOAD_ERROR');
		}
		if ($statusCode >= 300) {
			throw new Exception("DOWNLOAD_ERROR STATUS_CODE=$statusCode");
		}

		$fh = fopen($filePath, 'w');
		$writeResult = fwrite($fh, $res);
		if ($writeResult === false) {
			fclose($fh);
			throw new Exception('DOWNLOAD_FILE_WRITE_ERROR');
		}
		fclose($fh);
	}

	public function extract()
	{
		$basePath = $this->basePath;
		$zipPath = path_join($this->basePath, $this->downloadPath);

		$zip = new ZipArchive;
		$result = $zip->open($zipPath);
		if ($result !== true) {
			throw new Exception('ARCHIVE_FILE_OPEN_ERROR');
		}
		$result = $zip->extractTo($basePath);
		if ($result === false) {
			$zip->close();
			throw new Exception('ARCHIVE_EXTRACT_ERROR');
		}
		$zip->close();
	}

	public function place()
	{
		$basePath = $this->basePath;
		$extractedPath = path_join($this->basePath, $this->extractedDir);
		$zipPath = path_join($this->basePath, $this->downloadPath);

		$output = array();
		$code = 0;

		$cmd = "mv -f $extractedPath/* $basePath 2>/dev/null";
		exec($cmd, $output, $code);
		if ($code !== 0) {
			// For directory overwrite
			$cmd = "cp -Rf $extractedPath/* $basePath 2>/dev/null";
			exec($cmd, $output, $code);
		}
		if ($code !== 0) {
			throw new Exception('PLACE_ERROR(CMS_FILES_PLACE)');
		}

		$cmd = "rm -Rf $extractedPath 2>/dev/null";
		exec($cmd, $output, $code);
		if ($code !== 0) {
			throw new Exception('PLACE_ERROR(EMPTY_DIR_DELETE)');
		}

		$cmd = "rm -Rf $zipPath 2>/dev/null";
		exec($cmd, $output, $code);
		if ($code !== 0) {
			throw new Exception('PLACE_ERROR(ZIP_FILE_DELETE)');
		}
	}

	public function install()
	{
		$this->download();
		$this->extract();
		$this->place();
	}
}

class Authenticator
{
	public $token = TOKEN;
	public $expiredAt = EXPIRED_AT;

	public function exec($token, $now = null)
	{
		if (!$now) {
			$now = strtotime('now');
		}
		if ($this->token !== $token) {
			throw new Exception("AUTH_ERROR");
		}
		if ((int)$this->expiredAt <= (int)$now) {
			throw new Exception("AUTH_ERROR");
		}
		return;
	}
}

function main()
{
	ini_set('display_errors', 0);

	$res = array(
		'result' => true,
		'error_message' => ''
	);

	try {
		$token = $_POST['token'];
		$authenticator = new Authenticator();
		$authenticator->exec($token);
		$cmsArchive = new CmsArchive();
		$cmsArchive->install();
	} catch(Exception $e) {
		$res['result'] = false;
		$res['error_message'] = $e->getMessage();
	}

	echo json_encode($res);
}

if (isset($_POST['token'])) {
	main();
}
`
