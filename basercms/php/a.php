<?php
// 2>&1 をつけることで 標準エラー出力を標準出力に向ける
$cmd = "php b.php 2>&1";

$output = [];
$exit = 0;

// $outputには標準出力結果が入る。改行で要素分割して配列に追加される。
// $exitはexitcode
if (!exec($cmd, $output, $exit)) {
    // コマンドが存在しないときや、呼び出したスクリプトがエラーしたときは
    // ここが実行される
    // 故意にexitコードを投げても発火しない？
    echo "EXEC_ERROR";
}
?>

<pre>
<?php echo json_encode($output); ?>
</pre>

<hr>

<?php echo $exit; ?>

<hr>
