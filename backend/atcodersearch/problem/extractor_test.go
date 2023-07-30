package problem

import (
	"strings"
	"testing"
)

func TestExtract(t *testing.T) {
	source := `
<!DOCTYPE html>
<html>
<head>
	<title>C - Ideal Sheet</title>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta http-equiv="Content-Language" content="ja">
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
	<meta name="format-detection" content="telephone=no">
	<meta name="google-site-verification" content="nXGC_JxO0yoP1qBzMnYD_xgufO6leSLw1kyNo2HZltM" />

	
	<script async src="https://www.googletagmanager.com/gtag/js?id=G-RC512FD18N"></script>
	<script>
		window.dataLayer = window.dataLayer || [];
		function gtag(){dataLayer.push(arguments);}
		gtag('js', new Date());

		gtag('config', 'G-RC512FD18N');
	</script>

	
	<meta name="description" content="プログラミング初級者から上級者まで楽しめる、競技プログラミングコンテストサイト「AtCoder」。オンラインで毎週開催プログラミングコンテストを開催しています。競技プログラミングを用いて、客観的に自分のスキルを計ることのできるサービスです。">
	<meta name="author" content="AtCoder Inc.">

	<meta property="og:site_name" content="AtCoder">
	
	<meta property="og:title" content="C - Ideal Sheet" />
	<meta property="og:description" content="プログラミング初級者から上級者まで楽しめる、競技プログラミングコンテストサイト「AtCoder」。オンラインで毎週開催プログラミングコンテストを開催しています。競技プログラミングを用いて、客観的に自分のスキルを計ることのできるサービスです。" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://atcoder.jp/contests/abc307/tasks/abc307_c" />
	<meta property="og:image" content="https://img.atcoder.jp/assets/atcoder.png" />
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:site" content="@atcoder" />
	
	<meta property="twitter:title" content="C - Ideal Sheet" />

	<link href="//fonts.googleapis.com/css?family=Lato:400,700" rel="stylesheet" type="text/css">
	<link rel="stylesheet" type="text/css" href="//img.atcoder.jp/public/646bce9/css/bootstrap.min.css">
	<link rel="stylesheet" type="text/css" href="//img.atcoder.jp/public/646bce9/css/base.css">
	<link rel="shortcut icon" type="image/png" href="//img.atcoder.jp/assets/favicon.png">
	<link rel="apple-touch-icon" href="//img.atcoder.jp/assets/atcoder.png">
	<script src="//img.atcoder.jp/public/646bce9/js/lib/jquery-1.9.1.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/lib/bootstrap.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/cdn/js.cookie.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/cdn/moment.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/cdn/moment_js-ja.js"></script>
	<script>
		var LANG = "ja";
		var userScreenName = "fjnkt98";
		var csrfToken = "J0H5Y5m/RVVl3DcGJmCaMbXU1XlSlNSccxiajEiH9pg="
	</script>
	<script src="//img.atcoder.jp/public/646bce9/js/utils.js"></script>
	
	
		<script src="//img.atcoder.jp/public/646bce9/js/contest.js"></script>
		<link href="//img.atcoder.jp/public/646bce9/css/contest.css" rel="stylesheet" />
		<script>
			var contestScreenName = "abc307";
			var remainingText = "残り時間";
			var countDownText = "開始まであと";
			var startTime = moment("2023-06-24T21:00:00+09:00");
			var endTime = moment("2023-06-24T22:40:00+09:00");
		</script>
		<style></style>
	
	
		<link href="//img.atcoder.jp/public/646bce9/css/cdn/select2.min.css" rel="stylesheet" />
		<link href="//img.atcoder.jp/public/646bce9/css/cdn/select2-bootstrap.min.css" rel="stylesheet" />
		<script src="//img.atcoder.jp/public/646bce9/js/lib/select2.min.js"></script>
	
	
		<link rel="stylesheet" href="//img.atcoder.jp/public/646bce9/css/cdn/codemirror.min.css">
		<script src="//img.atcoder.jp/public/646bce9/js/cdn/codemirror.min.js"></script>
		<script src="//img.atcoder.jp/public/646bce9/js/codeMirror/merged.js"></script>
	
	
		<script src="//img.atcoder.jp/public/646bce9/js/cdn/run_prettify.js"></script>
	
	
		<link rel="stylesheet" href="//img.atcoder.jp/public/646bce9/css/cdn/katex.min.css">
		<script defer src="//img.atcoder.jp/public/646bce9/js/cdn/katex.min.js"></script>
		<script defer src="//img.atcoder.jp/public/646bce9/js/cdn/auto-render.min.js"></script>
		<script>$(function(){$('var').each(function(){var html=$(this).html().replace(/<sub>/g,'_{').replace(/<\/sub>/g,'}');$(this).html('\\('+html+'\\)');});});</script>
		<script>
			var katexOptions = {
				delimiters: [
					{left: "$$", right: "$$", display: true},
					
					{left: "\\(", right: "\\)", display: false},
					{left: "\\[", right: "\\]", display: true}
				],
      	ignoredTags: ["script", "noscript", "style", "textarea", "code", "option"],
				ignoredClasses: ["prettyprint", "source-code-for-copy"],
				throwOnError: false
			};
			document.addEventListener("DOMContentLoaded", function() { renderMathInElement(document.body, katexOptions);});
		</script>
	
	
	
	
	
	
	
	
	
	<script src="//img.atcoder.jp/public/646bce9/js/base.js"></script>
</head>

<body>

<script type="text/javascript">
	var __pParams = __pParams || [];
	__pParams.push({client_id: '468', c_1: 'atcodercontest', c_2: 'ClientSite'});
</script>
<script type="text/javascript" src="https://cdn.d2-apps.net/js/tr.js" async></script>


<div id="modal-contest-start" class="modal fade" tabindex="-1" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
				<h4 class="modal-title">コンテスト開始</h4>
			</div>
			<div class="modal-body">
				<p>東京海上日動プログラミングコンテスト2023（AtCoder Beginner Contest 307）が開始されました。</p>
			</div>
			<div class="modal-footer">
				
					<button type="button" class="btn btn-default" data-dismiss="modal">閉じる</button>
				
			</div>
		</div>
	</div>
</div>
<div id="modal-contest-end" class="modal fade" tabindex="-1" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
				<h4 class="modal-title">コンテスト終了</h4>
			</div>
			<div class="modal-body">
				<p>東京海上日動プログラミングコンテスト2023（AtCoder Beginner Contest 307）は終了しました。</p>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">閉じる</button>
			</div>
		</div>
	</div>
</div>
<div id="main-div" class="float-container">


	<nav class="navbar navbar-inverse navbar-fixed-top">
		<div class="container-fluid">
			<div class="navbar-header">
				<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar-collapse" aria-expanded="false">
					<span class="icon-bar"></span><span class="icon-bar"></span><span class="icon-bar"></span>
				</button>
				<a class="navbar-brand" href="/home"></a>
			</div>
			<div class="collapse navbar-collapse" id="navbar-collapse">
				<ul class="nav navbar-nav">
				
					<li><a class="contest-title" href="/contests/abc307">東京海上日動プログラミングコンテスト2023（AtCoder Beginner Contest 307）</a></li>
				
				</ul>
				<ul class="nav navbar-nav navbar-right">
					
					<li class="dropdown">
						<a class="dropdown-toggle" data-toggle="dropdown" href="#" role="button" aria-haspopup="true" aria-expanded="false">
							<img src='//img.atcoder.jp/assets/top/img/flag-lang/ja.png'> 日本語 <span class="caret"></span>
						</a>
						<ul class="dropdown-menu">
							<li><a href="/contests/abc307/tasks/abc307_c?lang=ja"><img src='//img.atcoder.jp/assets/top/img/flag-lang/ja.png'> 日本語</a></li>
							<li><a href="/contests/abc307/tasks/abc307_c?lang=en"><img src='//img.atcoder.jp/assets/top/img/flag-lang/en.png'> English</a></li>
						</ul>
					</li>
					
					
						<li class="dropdown">
							<a class="dropdown-toggle" data-toggle="dropdown" href="#" role="button" aria-haspopup="true" aria-expanded="false">
								<span class="glyphicon glyphicon-cog" aria-hidden="true"></span> fjnkt98 (Guest) <span class="caret"></span>
							</a>
							<ul class="dropdown-menu">
								<li><a href="/users/fjnkt98"><span class="glyphicon glyphicon-user" aria-hidden="true"></span> マイプロフィール</a></li>
								<li class="divider"></li>
								<li><a href="/settings"><span class="glyphicon glyphicon-wrench" aria-hidden="true"></span> 基本設定</a></li>
								<li><a href="/settings/icon"><span class="glyphicon glyphicon-picture" aria-hidden="true"></span> アイコン設定</a></li>
								<li><a href="/settings/password"><span class="glyphicon glyphicon-lock" aria-hidden="true"></span> パスワードの変更</a></li>
								<li><a href="/settings/fav"><span class="glyphicon glyphicon-star" aria-hidden="true"></span> お気に入り管理</a></li>
								
								
								
								<li class="divider"></li>
								<li><a href="javascript:void(form_logout.submit())"><span class="glyphicon glyphicon-log-out" aria-hidden="true"></span> ログアウト</a></li>
							</ul>
						</li>
					
				</ul>
			</div>
		</div>
	</nav>

	<form method="POST" name="form_logout" action="/logout?continue=https%3A%2F%2Fatcoder.jp%2Fcontests%2Fabc307%2Ftasks%2Fabc307_c">
		<input type="hidden" name="csrf_token" value="J0H5Y5m/RVVl3DcGJmCaMbXU1XlSlNSccxiajEiH9pg=" />
	</form>
	<div id="main-container" class="container"
		 	style="padding-top:50px;">
		


<div class="row">
	<div id="contest-nav-tabs" class="col-sm-12 mb-2 cnvtb-fixed">
	<div>
		<small class="contest-duration">
			
				コンテスト時間:
				<a href='http://www.timeanddate.com/worldclock/fixedtime.html?iso=20230624T2100&p1=248' target='blank'><time class='fixtime fixtime-full'>2023-06-24 21:00:00+0900</time></a> ~ <a href='http://www.timeanddate.com/worldclock/fixedtime.html?iso=20230624T2240&p1=248' target='blank'><time class='fixtime fixtime-full'>2023-06-24 22:40:00+0900</time></a> 
				(100分)
			
		</small>
		<small class="back-to-home pull-right"><a href="/home">AtCoderホームへ戻る</a></small>
	</div>
	<ul class="nav nav-tabs">
		<li><a href="/contests/abc307"><span class="glyphicon glyphicon-home" aria-hidden="true"></span> トップ</a></li>
		
			<li class="active"><a href="/contests/abc307/tasks"><span class="glyphicon glyphicon-tasks" aria-hidden="true"></span> 問題</a></li>
		

		
			<li><a href="/contests/abc307/clarifications"><span class="glyphicon glyphicon-question-sign" aria-hidden="true"></span> 質問 <span id="clar-badge" class="badge"></span></a></li>
		

		
			<li><a href="/contests/abc307/submit?taskScreenName=abc307_c"><span class="glyphicon glyphicon-send" aria-hidden="true"></span> 提出</a></li>
		

		
			<li>
				<a class="dropdown-toggle" data-toggle="dropdown" href="#" role="button" aria-haspopup="true" aria-expanded="false"><span class="glyphicon glyphicon-list" aria-hidden="true"></span> 提出結果<span class="caret"></span></a>
				<ul class="dropdown-menu">
					<li><a href="/contests/abc307/submissions"><span class="glyphicon glyphicon-globe" aria-hidden="true"></span> すべての提出</a></li>
					
						<li><a href="/contests/abc307/submissions/me"><span class="glyphicon glyphicon-user" aria-hidden="true"></span> 自分の提出</a></li>
						<li class="divider"></li>
						<li><a href="/contests/abc307/score"><span class="glyphicon glyphicon-dashboard" aria-hidden="true"></span> 自分の得点状況</a></li>
					
				</ul>
			</li>
		

		
			
				
					<li><a href="/contests/abc307/standings"><span class="glyphicon glyphicon-sort-by-attributes-alt" aria-hidden="true"></span> 順位表</a></li>
				
			
				
					<li><a href="/contests/abc307/standings/virtual"><span class="glyphicon glyphicon-sort-by-attributes-alt" aria-hidden="true"></span> バーチャル順位表</a></li>
				
			
		

		
			<li><a href="/contests/abc307/custom_test"><span class="glyphicon glyphicon-wrench" aria-hidden="true"></span> コードテスト</a></li>
		

		
			<li><a href="/contests/abc307/editorial"><span class="glyphicon glyphicon-book" aria-hidden="true"></span> 解説</a></li>
		
		

		<li class="pull-right"><a id="fix-cnvtb" href="javascript:void(0)"><span class="glyphicon glyphicon-pushpin" aria-hidden="true"></span></a></li>
	</ul>
</div>
	<div class="col-sm-12">
		<span class="h2">
			C - Ideal Sheet
			<a class="btn btn-default btn-sm" href="/contests/abc307/tasks/abc307_c/editorial">解説</a>
		</span>
		<span id="task-lang-btn" class="pull-right"><span data-lang="ja"><img src='//img.atcoder.jp/assets/top/img/flag-lang/ja.png'></span> / <span data-lang="en"><img src='//img.atcoder.jp/assets/top/img/flag-lang/en.png'></span></span>
		<script>
			$(function() {
				var ts = $('#task-statement span.lang');
				if (ts.children('span').size() <= 1) {
					$('#task-lang-btn').hide();
					ts.children('span').show();
					return;
				}
				var REMEMBER_LB = 5;
				var LS_KEY = 'task_lang';
				var taskLang = getLS(LS_KEY) || '';
				var changeTimes = 0;
				if (taskLang == 'ja' || taskLang == 'en') {
					changeTimes = REMEMBER_LB;
				} else {
					var changeTimes = parseInt(taskLang, 10);
					if (isNaN(changeTimes)) {
						changeTimes = 0;
						delLS(LS_KEY);
					}
					changeTimes++;
					taskLang = LANG;
				}
				ts.children('span.lang-' + taskLang).show();

				$('#task-lang-btn span').click(function() {
					var l = $(this).data('lang');
					ts.children('span').hide();
					ts.children('span.lang-' + l).show();
					if (changeTimes < REMEMBER_LB) setLS(LS_KEY, changeTimes);
					else setLS(LS_KEY, l);
				});
			});
		</script>
		<hr/>
		<p>
			実行時間制限: 2 sec / メモリ制限: 1024 MB
			
		</p>

		<div id="task-statement">
			
			<span class="lang">
<span class="lang-ja">
<p>配点 : <var>300</var> 点</p>

<div class="part">
<section>
<h3>問題文</h3><p>高橋君は黒いマスと透明なマスからなるシート <var>A,B</var> を <var>1</var> 枚ずつと、透明なマスのみからなる無限に広がるシート <var>C</var> を持っています。<br />
また、高橋君には黒いマスと透明なマスからなる、理想とするシート <var>X</var> が存在します。</p>
<p>シート <var>A,B,X</var> の大きさはそれぞれ縦 <var>H_A</var> マス <var>\times</var> 横 <var>W_A</var> マス、縦 <var>H_B</var> マス <var>\times</var> 横 <var>W_B</var> マス、縦 <var>H_X</var> マス <var>\times</var> 横 <var>W_X</var> マスです。<br />
シート <var>A</var> の各マスは <code>.</code> と <code>#</code> からなる長さ <var>W_A</var> の文字列 <var>H_A</var> 個 <var>A_1,A_2,\ldots,A_{H_A}</var> によって表され、<br />
<var>A_i</var> <var>(1\leq i\leq H_A)</var> の <var>j</var> 文字目 <var>(1\leq j\leq W_A)</var> が、
 <code>.</code> のときシート <var>A</var> の上から <var>i</var> 行目かつ左から <var>j</var> 列目のマスは透明なマスであり、
<code>#</code> のとき黒いマスです。<br />
シート <var>B,X</var> の各マスも、同様に長さ <var>W_B</var> の文字列 <var>H_B</var> 個 <var>B_1,B_2,\ldots,B_{H_B}</var> および長さ <var>W_X</var> の文字列 <var>H_X</var> 個 <var>X_1,X_2,\ldots,X_{H_X}</var> によって表されます。</p>
<p>高橋君の目標は、次の手順で、シート <var>A,B,C</var> から、<var>A,B</var> に存在する <strong>すべての黒いマスを使って</strong> シート <var>X</var> を作り出すことです。</p>
<ol>
<li>シート <var>A,B</var> をマス目に沿ってシート <var>C</var> に貼り付ける。この時、シート <var>A,B</var> はそれぞれ好きな場所に平行移動させて貼って良いが、シートを切り分けたり、回転させたりしてはいけない。</li>
<li>シート <var>C</var> からマス目に沿って <var>H_X\times W_X</var> マスの領域を切り出す。ここで、切り出されたシートの各マスは、シート <var>A</var> または <var>B</var> の黒いマスが貼り付けられていれば黒いマスに、そうでなければ透明なマスとなる。</li>
</ol>
<p>このとき、貼り付ける位置と切り出す領域をうまくとることで高橋君は目標を達成できるか、すなわち次の条件をともにみたすことにできるか判定してください。</p>
<ul>
<li>切り出されたシートはシート <var>A,B</var> の <strong>黒いマスをすべて</strong> 含む。切り出されたシートの上でシート <var>A,B</var> の黒いマスどうしが重なって存在していても構わない。</li>
<li>切り出されたシートは、回転させたり裏返したりすることなくシート <var>X</var> と一致する。</li>
</ul>
</section>
</div>

<div class="part">
<section>
<h3>制約</h3><ul>
<li><var>1\leq H_A,W_A,H_B,W_B,H_X,W_X\leq 10</var></li>
<li><var>H_A,W_A,H_B,W_B,H_X,W_X</var> は整数</li>
<li><var>A_i</var> は <code>.</code> と <code>#</code> のみからなる長さ <var>W_A</var> の文字列</li>
<li><var>B_i</var> は <code>.</code> と <code>#</code> のみからなる長さ <var>W_B</var> の文字列</li>
<li><var>X_i</var> は <code>.</code> と <code>#</code> のみからなる長さ <var>W_X</var> の文字列</li>
<li>シート <var>A,B,X</var> はそれぞれ少なくとも <var>1</var> つ以上の黒いマスを含む。</li>
</ul>
</section>
</div>

<hr />
<div class="io-style">
<div class="part">
<section>
<h3>入力</h3><p>入力は以下の形式で標準入力から与えられる。</p>
<pre><var>H_A</var> <var>W_A</var>
<var>A_1</var>
<var>A_2</var>
<var>\vdots</var>
<var>A_{H_A}</var>
<var>H_B</var> <var>W_B</var>
<var>B_1</var>
<var>B_2</var>
<var>\vdots</var>
<var>B_{H_B}</var>
<var>H_X</var> <var>W_X</var>
<var>X_1</var>
<var>X_2</var>
<var>\vdots</var>
<var>X_{H_X}</var>
</pre>
</section>
</div>

<div class="part">
<section>
<h3>出力</h3><p>高橋君が問題文中の目標を達成できるならば <code>Yes</code> を、できないならば <code>No</code> を出力せよ。</p>
</section>
</div>
</div>

<hr />
<div class="part">
<section>
<h3>入力例 1</h3><pre>3 5
#.#..
.....
.#...
2 2
#.
.#
5 3
...
#.#
.#.
.#.
...
</pre>
</section>
</div>

<div class="part">
<section>
<h3>出力例 1</h3><pre>Yes
</pre>
<p>まず、シート <var>A</var> をシート <var>C</var> に貼り付けると下図のようになります。</p>
<pre>     <var>\vdots</var>
  .......  
  .#.#...  
<var>\cdots</var>.......<var>\cdots</var>
  ..#....  
  .......  
     <var>\vdots</var>
</pre>
<p>さらに、シート <var>B</var> をシート <var>A</var> と左上を合わせて貼ってみると下図のようになります。</p>
<pre>     <var>\vdots</var>
  .......  
  .#.#...  
<var>\cdots</var>..#....<var>\cdots</var>
  ..#....  
  .......  
     <var>\vdots</var>
</pre>
<p>ここで、上で具体的に図示されている範囲のうち、上から <var>1</var> 行目かつ左から <var>2</var> 列目のマスを左上として
<var>5\times 3</var> マスを切り出すと下図のようになります。</p>
<pre>...
#.#
.#.
.#.
...
</pre>
<p>これはシート <var>A,B</var> のすべての黒いマスを含んでおり、また、シート <var>X</var> と一致しているため条件を満たしています。</p>
<p>よって、<code>Yes</code> を出力します。</p>
</section>
</div>

<hr />
<div class="part">
<section>
<h3>入力例 2</h3><pre>2 2
#.
.#
2 2
#.
.#
2 2
##
##
</pre>
</section>
</div>

<div class="part">
<section>
<h3>出力例 2</h3><pre>No
</pre>
<p>シート <var>A</var> や <var>B</var> を回転させて貼ってはいけないことに注意してください。</p>
</section>
</div>

<hr />
<div class="part">
<section>
<h3>入力例 3</h3><pre>1 1
#
1 2
##
1 1
#
</pre>
</section>
</div>

<div class="part">
<section>
<h3>出力例 3</h3><pre>No
</pre>
<p>どのように貼ったり切り出したりしても、シート <var>B</var> の黒いマスをすべて含むように切り出すことはできないため、<var>1</var> つめの条件をみたすことができません。
よって、<code>No</code> を出力します。</p>
</section>
</div>

<hr />
<div class="part">
<section>
<h3>入力例 4</h3><pre>3 3
###
...
...
3 3
#..
#..
#..
3 3
..#
..#
###
</pre>
</section>
</div>

<div class="part">
<section>
<h3>出力例 4</h3><pre>Yes
</pre></section>
</div>
</span>
<span class="lang-en">
<p>Score : <var>300</var> points</p>

<div class="part">
<section>
<h3>Problem Statement</h3><p>Takahashi has two sheets <var>A</var> and <var>B</var>, each composed of black squares and transparent squares, and an infinitely large sheet <var>C</var> composed of transparent squares.<br />
There is also an ideal sheet <var>X</var> for Takahashi composed of black squares and transparent squares.</p>
<p>The sizes of sheets <var>A</var>, <var>B</var>, and <var>X</var> are <var>H_A</var> rows <var>\times</var> <var>W_A</var> columns, <var>H_B</var> rows <var>\times</var> <var>W_B</var> columns, and <var>H_X</var> rows <var>\times</var> <var>W_X</var> columns, respectively.<br />
The squares of sheet <var>A</var> are represented by <var>H_A</var> strings of length <var>W_A</var>, <var>A_1, A_2, \ldots, A_{H_A}</var> consisting of <code>.</code> and <code>#</code>.<br />
If the <var>j</var>-th character <var>(1\leq j\leq W_A)</var> of <var>A_i</var> <var>(1\leq i\leq H_A)</var> is <code>.</code>, the square at the <var>i</var>-th row from the top and <var>j</var>-th column from the left is transparent; if it is <code>#</code>, that square is black.<br />
Similarly, the squares of sheets <var>B</var> and <var>X</var> are represented by <var>H_B</var> strings of length <var>W_B</var>, <var>B_1, B_2, \ldots, B_{H_B}</var>, and <var>H_X</var> strings of length <var>W_X</var>, <var>X_1, X_2, \ldots, X_{H_X}</var>, respectively.</p>
<p>Takahashi's goal is to create sheet <var>X</var> using <strong>all black squares</strong> in sheets <var>A</var> and <var>B</var> by following the steps below with sheets <var>A</var>, <var>B</var>, and <var>C</var>.</p>
<ol>
<li>Paste sheets <var>A</var> and <var>B</var> onto sheet <var>C</var> along the grid. Each sheet can be pasted anywhere by translating it, but it cannot be cut or rotated.</li>
<li>Cut out an <var>H_X\times W_X</var> area from sheet <var>C</var> along the grid. Here, a square of the cut-out sheet will be black if a black square of sheet <var>A</var> or <var>B</var> is pasted there, and transparent otherwise.</li>
</ol>
<p>Determine whether Takahashi can achieve his goal by appropriately choosing the positions where the sheets are pasted and the area to cut out, that is, whether he can satisfy both of the following conditions.</p>
<ul>
<li>The cut-out sheet includes <strong>all black squares</strong> of sheets <var>A</var> and <var>B</var>. The black squares of sheets <var>A</var> and <var>B</var> may overlap on the cut-out sheet.</li>
<li>The cut-out sheet coincides sheet <var>X</var> without rotating or flipping.</li>
</ul>
</section>
</div>

<div class="part">
<section>
<h3>Constraints</h3><ul>
<li><var>1\leq H_A, W_A, H_B, W_B, H_X, W_X\leq 10</var></li>
<li><var>H_A, W_A, H_B, W_B, H_X, W_X</var> are integers.</li>
<li><var>A_i</var> is a string of length <var>W_A</var> consisting of <code>.</code> and <code>#</code>.</li>
<li><var>B_i</var> is a string of length <var>W_B</var> consisting of <code>.</code> and <code>#</code>.</li>
<li><var>X_i</var> is a string of length <var>W_X</var> consisting of <code>.</code> and <code>#</code>.</li>
<li>Sheets <var>A</var>, <var>B</var>, and <var>X</var> each contain at least one black square.</li>
</ul>
</section>
</div>

<hr />
<div class="io-style">
<div class="part">
<section>
<h3>Input</h3><p>The input is given from Standard Input in the following format:</p>
<pre><var>H_A</var> <var>W_A</var>
<var>A_1</var>
<var>A_2</var>
<var>\vdots</var>
<var>A_{H_A}</var>
<var>H_B</var> <var>W_B</var>
<var>B_1</var>
<var>B_2</var>
<var>\vdots</var>
<var>B_{H_B}</var>
<var>H_X</var> <var>W_X</var>
<var>X_1</var>
<var>X_2</var>
<var>\vdots</var>
<var>X_{H_X}</var>
</pre>
</section>
</div>

<div class="part">
<section>
<h3>Output</h3><p>If Takahashi can achieve the goal described in the problem statement, print <code>Yes</code>; otherwise, print <code>No</code>.</p>
</section>
</div>
</div>

<hr />
<div class="part">
<section>
<h3>Sample Input 1</h3><pre>3 5
#.#..
.....
.#...
2 2
#.
.#
5 3
...
#.#
.#.
.#.
...
</pre>
</section>
</div>

<div class="part">
<section>
<h3>Sample Output 1</h3><pre>Yes
</pre>
<p>First, paste sheet <var>A</var> onto sheet <var>C</var>, as shown in the figure below.</p>
<pre>     <var>\vdots</var>
  .......  
  .#.#...  
<var>\cdots</var>.......<var>\cdots</var>
  ..#....  
  .......  
     <var>\vdots</var>
</pre>
<p>Next, paste sheet <var>B</var> so that its top-left corner aligns with that of sheet <var>A</var>, as shown in the figure below.</p>
<pre>     <var>\vdots</var>
  .......  
  .#.#...  
<var>\cdots</var>..#....<var>\cdots</var>
  ..#....  
  .......  
     <var>\vdots</var>
</pre>
<p>Now, cut out a <var>5\times 3</var> area with the square in the first row and second column of the range illustrated above as the top-left corner, as shown in the figure below.</p>
<pre>...
#.#
.#.
.#.
...
</pre>
<p>This includes all black squares of sheets <var>A</var> and <var>B</var> and matches sheet <var>X</var>, satisfying the conditions.</p>
<p>Therefore, print <code>Yes</code>.</p>
</section>
</div>

<hr />
<div class="part">
<section>
<h3>Sample Input 2</h3><pre>2 2
#.
.#
2 2
#.
.#
2 2
##
##
</pre>
</section>
</div>

<div class="part">
<section>
<h3>Sample Output 2</h3><pre>No
</pre>
<p>Note that sheets <var>A</var> and <var>B</var> may not be rotated or flipped when pasting them.</p>
</section>
</div>

<hr />
<div class="part">
<section>
<h3>Sample Input 3</h3><pre>1 1
#
1 2
##
1 1
#
</pre>
</section>
</div>

<div class="part">
<section>
<h3>Sample Output 3</h3><pre>No
</pre>
<p>No matter how you paste or cut, you cannot cut out a sheet that includes all black squares of sheet <var>B</var>, so you cannot satisfy the first condition.
Therefore, print <code>No</code>.</p>
</section>
</div>

<hr />
<div class="part">
<section>
<h3>Sample Input 4</h3><pre>3 3
###
...
...
3 3
#..
#..
#..
3 3
..#
..#
###
</pre>
</section>
</div>

<div class="part">
<section>
<h3>Sample Output 4</h3><pre>Yes
</pre></section>
</div>
</span>
</span>

		</div>

		

		
		<hr/>
		<form class="form-horizontal form-code-submit" action="/contests/abc307/submit" method="POST">
			<input type="hidden" name="data.TaskScreenName" value="abc307_c" />
			
			<div class="form-group ">
				<label class="control-label col-sm-2" for="select-lang">言語</label>
				<div id="select-lang" class="col-sm-5" data-name="data.LanguageId">
					<select class="form-control current" data-placeholder="-" name="data.LanguageId" required>
						<option></option>
						
							<option value="4001" data-mime="text/x-csrc">C (GCC 9.2.1)</option>
						
							<option value="4002" data-mime="text/x-csrc">C (Clang 10.0.0)</option>
						
							<option value="4003" data-mime="text/x-c&#43;&#43;src">C&#43;&#43; (GCC 9.2.1)</option>
						
							<option value="4004" data-mime="text/x-c&#43;&#43;src">C&#43;&#43; (Clang 10.0.0)</option>
						
							<option value="4005" data-mime="text/x-java">Java (OpenJDK 11.0.6)</option>
						
							<option value="4006" data-mime="text/x-python">Python (3.8.2)</option>
						
							<option value="4007" data-mime="text/x-sh">Bash (5.0.11)</option>
						
							<option value="4008" data-mime="text/x-bc">bc (1.07.1)</option>
						
							<option value="4009" data-mime="text/x-sh">Awk (GNU Awk 4.1.4)</option>
						
							<option value="4010" data-mime="text/x-csharp">C# (.NET Core 3.1.201)</option>
						
							<option value="4011" data-mime="text/x-csharp">C# (Mono-mcs 6.8.0.105)</option>
						
							<option value="4012" data-mime="text/x-csharp">C# (Mono-csc 3.5.0)</option>
						
							<option value="4013" data-mime="text/x-clojure">Clojure (1.10.1.536)</option>
						
							<option value="4014" data-mime="text/x-crystal">Crystal (0.33.0)</option>
						
							<option value="4015" data-mime="text/x-d">D (DMD 2.091.0)</option>
						
							<option value="4016" data-mime="text/x-d">D (GDC 9.2.1)</option>
						
							<option value="4017" data-mime="text/x-d">D (LDC 1.20.1)</option>
						
							<option value="4018" data-mime="application/dart">Dart (2.7.2)</option>
						
							<option value="4019" data-mime="text/x-dc">dc (1.4.1)</option>
						
							<option value="4020" data-mime="text/x-erlang">Erlang (22.3)</option>
						
							<option value="4021" data-mime="elixir">Elixir (1.10.2)</option>
						
							<option value="4022" data-mime="text/x-fsharp">F# (.NET Core 3.1.201)</option>
						
							<option value="4023" data-mime="text/x-fsharp">F# (Mono 10.2.3)</option>
						
							<option value="4024" data-mime="text/x-forth">Forth (gforth 0.7.3)</option>
						
							<option value="4025" data-mime="text/x-fortran">Fortran (GNU Fortran 9.2.1)</option>
						
							<option value="4026" data-mime="text/x-go">Go (1.14.1)</option>
						
							<option value="4027" data-mime="text/x-haskell">Haskell (GHC 8.8.3)</option>
						
							<option value="4028" data-mime="text/x-haxe">Haxe (4.0.3); js</option>
						
							<option value="4029" data-mime="text/x-haxe">Haxe (4.0.3); Java</option>
						
							<option value="4030" data-mime="text/javascript">JavaScript (Node.js 12.16.1)</option>
						
							<option value="4031" data-mime="text/x-julia">Julia (1.4.0)</option>
						
							<option value="4032" data-mime="text/x-kotlin">Kotlin (1.3.71)</option>
						
							<option value="4033" data-mime="text/x-lua">Lua (Lua 5.3.5)</option>
						
							<option value="4034" data-mime="text/x-lua">Lua (LuaJIT 2.1.0)</option>
						
							<option value="4035" data-mime="text/x-sh">Dash (0.5.8)</option>
						
							<option value="4036" data-mime="text/x-nim">Nim (1.0.6)</option>
						
							<option value="4037" data-mime="text/x-objectivec">Objective-C (Clang 10.0.0)</option>
						
							<option value="4038" data-mime="text/x-common-lisp">Common Lisp (SBCL 2.0.3)</option>
						
							<option value="4039" data-mime="text/x-ocaml">OCaml (4.10.0)</option>
						
							<option value="4040" data-mime="text/x-octave">Octave (5.2.0)</option>
						
							<option value="4041" data-mime="text/x-pascal">Pascal (FPC 3.0.4)</option>
						
							<option value="4042" data-mime="text/x-perl">Perl (5.26.1)</option>
						
							<option value="4043" data-mime="text/x-perl">Raku (Rakudo 2020.02.1)</option>
						
							<option value="4044" data-mime="text/x-php">PHP (7.4.4)</option>
						
							<option value="4045" data-mime="text/x-prolog">Prolog (SWI-Prolog 8.0.3)</option>
						
							<option value="4046" data-mime="text/x-python">PyPy2 (7.3.0)</option>
						
							<option value="4047" data-mime="text/x-python">PyPy3 (7.3.0)</option>
						
							<option value="4048" data-mime="text/x-racket">Racket (7.6)</option>
						
							<option value="4049" data-mime="text/x-ruby">Ruby (2.7.1)</option>
						
							<option value="4050" data-mime="text/x-rustsrc">Rust (1.42.0)</option>
						
							<option value="4051" data-mime="text/x-scala">Scala (2.13.1)</option>
						
							<option value="4052" data-mime="text/x-java">Java (OpenJDK 1.8.0)</option>
						
							<option value="4053" data-mime="text/x-scheme">Scheme (Gauche 0.9.9)</option>
						
							<option value="4054" data-mime="text/x-sml">Standard ML (MLton 20130715)</option>
						
							<option value="4055" data-mime="text/x-swift">Swift (5.2.1)</option>
						
							<option value="4056" data-mime="text/plain">Text (cat 8.28)</option>
						
							<option value="4057" data-mime="text/typescript">TypeScript (3.8)</option>
						
							<option value="4058" data-mime="text/x-vb">Visual Basic (.NET Core 3.1.101)</option>
						
							<option value="4059" data-mime="text/x-sh">Zsh (5.4.2)</option>
						
							<option value="4060" data-mime="text/x-cobol">COBOL - Fixed (OpenCOBOL 1.1.0)</option>
						
							<option value="4061" data-mime="text/x-cobol">COBOL - Free (OpenCOBOL 1.1.0)</option>
						
							<option value="4062" data-mime="text/x-brainfuck">Brainfuck (bf 20041219)</option>
						
							<option value="4063" data-mime="text/x-ada">Ada2012 (GNAT 9.2.1)</option>
						
							<option value="4064" data-mime="text/x-unlambda">Unlambda (2.0.0)</option>
						
							<option value="4065" data-mime="text/x-python">Cython (0.29.16)</option>
						
							<option value="4066" data-mime="text/x-sh">Sed (4.4)</option>
						
							<option value="4067" data-mime="text/x-vim">Vim (8.2.0460)</option>
						
					</select>
					<span class="error"></span>
				</div>
			</div>
			<script>var currentLang = getLS('defaultLang');</script>
			
			
<div class="form-group">
	<label class="control-label col-sm-2" for="sourceCode">ソースコード</label>
	<div class="col-sm-7" id="sourceCode">
		<div class="div-editor">
			<textarea class="form-control editor" name="sourceCode"></textarea>
		</div>
		<textarea class="form-control plain-textarea" style="display:none;"></textarea>
		<p>
			<span class="gray">※ 512 KiB まで</span><br>
			<span class="gray">※ ソースコードは「Main.<i>拡張子</i>」で保存されます</span>
		</p>
	</div>
	<div class="col-sm-3 editor-buttons">
		<p><button id="btn-open-file" type="button" class="btn btn-default btn-sm">
			<span class="glyphicon glyphicon-folder-open" aria-hidden="true"></span> &nbsp; ファイルを開く
		</button></p>
		<p><button type="button" class="btn btn-default btn-sm btn-toggle-editor" data-toggle="button" aria-pressed="false" autocomplete="off">
			エディタ切り替え
		</button></p>
		<p><button type="button" class="btn btn-default btn-sm btn-auto-height" data-toggle="button" aria-pressed="false" autocomplete="off">
			高さ自動調節
		</button></p>
	</div>
	<input id="input-open-file" type="file" style="display:none;">
</div>

			<input type="hidden" name="csrf_token" value="J0H5Y5m/RVVl3DcGJmCaMbXU1XlSlNSccxiajEiH9pg=" />
			<div class="form-group">
				<label class="control-label col-sm-2" for="submit"></label>
				<div class="col-sm-5">
					<button type="submit" class="btn btn-primary" id="submit">提出</button>
				</div>
			</div>
		</form>
		
	</div>
</div>




		
			<hr>
			
			
			
<div class="a2a_kit a2a_kit_size_20 a2a_default_style pull-right" data-a2a-url="https://atcoder.jp/contests/abc307/tasks/abc307_c?lang=ja" data-a2a-title="C - Ideal Sheet">
	<a class="a2a_button_facebook"></a>
	<a class="a2a_button_twitter"></a>
	
		<a class="a2a_button_hatena"></a>
	
	<a class="a2a_dd" href="https://www.addtoany.com/share"></a>
</div>

		
		<script async src="//static.addtoany.com/menu/page.js"></script>
		
	</div> 
	<hr>
</div> 

	<div class="container" style="margin-bottom: 80px;">
			<footer class="footer">
			
				<ul>
					<li><a href="/contests/abc307/rules">ルール</a></li>
					<li><a href="/contests/abc307/glossary">用語集</a></li>
					
				</ul>
			
			<ul>
				<li><a href="/tos">利用規約</a></li>
				<li><a href="/privacy">プライバシーポリシー</a></li>
				<li><a href="/personal">個人情報保護方針</a></li>
				<li><a href="/company">企業情報</a></li>
				<li><a href="/faq">よくある質問</a></li>
				<li><a href="/contact">お問い合わせ</a></li>
				<li><a href="/documents/request">資料請求</a></li>
			</ul>
			<div class="text-center">
					<small id="copyright">Copyright Since 2012 &copy;<a href="http://atcoder.co.jp">AtCoder Inc.</a> All rights reserved.</small>
			</div>
			</footer>
	</div>
	<p id="fixed-server-timer" class="contest-timer"></p>
	<div id="scroll-page-top" style="display:none;"><span class="glyphicon glyphicon-arrow-up" aria-hidden="true"></span> ページトップ</div>

</body>
</html>
`

	extractor := NewFullTextExtractor()
	ja, en, err := extractor.Extract(strings.NewReader(source))
	if err != nil {
		t.Errorf("failed to extract statement :%s", err.Error())
		return
	}

	t.Errorf("ja: %v, len(ja): %d, en: %v, len(en): %d", ja, len(ja), en, len(en))
}

func TestExtractLegacyProblem(t *testing.T) {
	source := `









<!DOCTYPE html>
<html>
<head>
	<title>C - 風力観測</title>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta http-equiv="Content-Language" content="ja">
	<meta name="viewport" content="width=device-width,initial-scale=1.0">
	<meta name="format-detection" content="telephone=no">
	<meta name="google-site-verification" content="nXGC_JxO0yoP1qBzMnYD_xgufO6leSLw1kyNo2HZltM" />

	
	<script async src="https://www.googletagmanager.com/gtag/js?id=G-RC512FD18N"></script>
	<script>
		window.dataLayer = window.dataLayer || [];
		function gtag(){dataLayer.push(arguments);}
		gtag('js', new Date());

		gtag('config', 'G-RC512FD18N');
	</script>

	
	<meta name="description" content="プログラミング初級者から上級者まで楽しめる、競技プログラミングコンテストサイト「AtCoder」。オンラインで毎週開催プログラミングコンテストを開催しています。競技プログラミングを用いて、客観的に自分のスキルを計ることのできるサービスです。">
	<meta name="author" content="AtCoder Inc.">

	<meta property="og:site_name" content="AtCoder">
	
	<meta property="og:title" content="C - 風力観測" />
	<meta property="og:description" content="プログラミング初級者から上級者まで楽しめる、競技プログラミングコンテストサイト「AtCoder」。オンラインで毎週開催プログラミングコンテストを開催しています。競技プログラミングを用いて、客観的に自分のスキルを計ることのできるサービスです。" />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://atcoder.jp/contests/abc001/tasks/abc001_3" />
	<meta property="og:image" content="https://img.atcoder.jp/assets/atcoder.png" />
	<meta name="twitter:card" content="summary" />
	<meta name="twitter:site" content="@atcoder" />
	
	<meta property="twitter:title" content="C - 風力観測" />

	<link href="//fonts.googleapis.com/css?family=Lato:400,700" rel="stylesheet" type="text/css">
	<link rel="stylesheet" type="text/css" href="//img.atcoder.jp/public/646bce9/css/bootstrap.min.css">
	<link rel="stylesheet" type="text/css" href="//img.atcoder.jp/public/646bce9/css/base.css">
	<link rel="shortcut icon" type="image/png" href="//img.atcoder.jp/assets/favicon.png">
	<link rel="apple-touch-icon" href="//img.atcoder.jp/assets/atcoder.png">
	<script src="//img.atcoder.jp/public/646bce9/js/lib/jquery-1.9.1.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/lib/bootstrap.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/cdn/js.cookie.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/cdn/moment.min.js"></script>
	<script src="//img.atcoder.jp/public/646bce9/js/cdn/moment_js-ja.js"></script>
	<script>
		var LANG = "ja";
		var userScreenName = "fjnkt98";
		var csrfToken = "J0H5Y5m/RVVl3DcGJmCaMbXU1XlSlNSccxiajEiH9pg="
	</script>
	<script src="//img.atcoder.jp/public/646bce9/js/utils.js"></script>
	
	
		<script src="//img.atcoder.jp/public/646bce9/js/contest.js"></script>
		<link href="//img.atcoder.jp/public/646bce9/css/contest.css" rel="stylesheet" />
		<script>
			var contestScreenName = "abc001";
			var remainingText = "残り時間";
			var countDownText = "開始まであと";
			var startTime = moment("2013-10-12T21:00:00+09:00");
			var endTime = moment("2013-10-12T23:00:00+09:00");
		</script>
		<style></style>
	
	
		<link href="//img.atcoder.jp/public/646bce9/css/cdn/select2.min.css" rel="stylesheet" />
		<link href="//img.atcoder.jp/public/646bce9/css/cdn/select2-bootstrap.min.css" rel="stylesheet" />
		<script src="//img.atcoder.jp/public/646bce9/js/lib/select2.min.js"></script>
	
	
		<link rel="stylesheet" href="//img.atcoder.jp/public/646bce9/css/cdn/codemirror.min.css">
		<script src="//img.atcoder.jp/public/646bce9/js/cdn/codemirror.min.js"></script>
		<script src="//img.atcoder.jp/public/646bce9/js/codeMirror/merged.js"></script>
	
	
		<script src="//img.atcoder.jp/public/646bce9/js/cdn/run_prettify.js"></script>
	
	
		<link rel="stylesheet" href="//img.atcoder.jp/public/646bce9/css/cdn/katex.min.css">
		<script defer src="//img.atcoder.jp/public/646bce9/js/cdn/katex.min.js"></script>
		<script defer src="//img.atcoder.jp/public/646bce9/js/cdn/auto-render.min.js"></script>
		<script>$(function(){$('var').each(function(){var html=$(this).html().replace(/<sub>/g,'_{').replace(/<\/sub>/g,'}');$(this).html('\\('+html+'\\)');});});</script>
		<script>
			var katexOptions = {
				delimiters: [
					{left: "$$", right: "$$", display: true},
					
					{left: "\\(", right: "\\)", display: false},
					{left: "\\[", right: "\\]", display: true}
				],
      	ignoredTags: ["script", "noscript", "style", "textarea", "code", "option"],
				ignoredClasses: ["prettyprint", "source-code-for-copy"],
				throwOnError: false
			};
			document.addEventListener("DOMContentLoaded", function() { renderMathInElement(document.body, katexOptions);});
		</script>
	
	
	
	
	
	
	
	
	
	<script src="//img.atcoder.jp/public/646bce9/js/base.js"></script>
</head>

<body>

<script type="text/javascript">
	var __pParams = __pParams || [];
	__pParams.push({client_id: '468', c_1: 'atcodercontest', c_2: 'ClientSite'});
</script>
<script type="text/javascript" src="https://cdn.d2-apps.net/js/tr.js" async></script>


<div id="modal-contest-start" class="modal fade" tabindex="-1" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
				<h4 class="modal-title">コンテスト開始</h4>
			</div>
			<div class="modal-body">
				<p>AtCoder Beginner Contest 001が開始されました。</p>
			</div>
			<div class="modal-footer">
				
					<button type="button" class="btn btn-default" data-dismiss="modal">閉じる</button>
				
			</div>
		</div>
	</div>
</div>
<div id="modal-contest-end" class="modal fade" tabindex="-1" role="dialog">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
				<h4 class="modal-title">コンテスト終了</h4>
			</div>
			<div class="modal-body">
				<p>AtCoder Beginner Contest 001は終了しました。</p>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">閉じる</button>
			</div>
		</div>
	</div>
</div>
<div id="main-div" class="float-container">


	<nav class="navbar navbar-inverse navbar-fixed-top">
		<div class="container-fluid">
			<div class="navbar-header">
				<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar-collapse" aria-expanded="false">
					<span class="icon-bar"></span><span class="icon-bar"></span><span class="icon-bar"></span>
				</button>
				<a class="navbar-brand" href="/home"></a>
			</div>
			<div class="collapse navbar-collapse" id="navbar-collapse">
				<ul class="nav navbar-nav">
				
					<li><a class="contest-title" href="/contests/abc001">AtCoder Beginner Contest 001</a></li>
				
				</ul>
				<ul class="nav navbar-nav navbar-right">
					
					<li class="dropdown">
						<a class="dropdown-toggle" data-toggle="dropdown" href="#" role="button" aria-haspopup="true" aria-expanded="false">
							<img src='//img.atcoder.jp/assets/top/img/flag-lang/ja.png'> 日本語 <span class="caret"></span>
						</a>
						<ul class="dropdown-menu">
							<li><a href="/contests/abc001/tasks/abc001_3?lang=ja"><img src='//img.atcoder.jp/assets/top/img/flag-lang/ja.png'> 日本語</a></li>
							<li><a href="/contests/abc001/tasks/abc001_3?lang=en"><img src='//img.atcoder.jp/assets/top/img/flag-lang/en.png'> English</a></li>
						</ul>
					</li>
					
					
						<li class="dropdown">
							<a class="dropdown-toggle" data-toggle="dropdown" href="#" role="button" aria-haspopup="true" aria-expanded="false">
								<span class="glyphicon glyphicon-cog" aria-hidden="true"></span> fjnkt98 (Guest) <span class="caret"></span>
							</a>
							<ul class="dropdown-menu">
								<li><a href="/users/fjnkt98"><span class="glyphicon glyphicon-user" aria-hidden="true"></span> マイプロフィール</a></li>
								<li class="divider"></li>
								<li><a href="/settings"><span class="glyphicon glyphicon-wrench" aria-hidden="true"></span> 基本設定</a></li>
								<li><a href="/settings/icon"><span class="glyphicon glyphicon-picture" aria-hidden="true"></span> アイコン設定</a></li>
								<li><a href="/settings/password"><span class="glyphicon glyphicon-lock" aria-hidden="true"></span> パスワードの変更</a></li>
								<li><a href="/settings/fav"><span class="glyphicon glyphicon-star" aria-hidden="true"></span> お気に入り管理</a></li>
								
								
								
								<li class="divider"></li>
								<li><a href="javascript:void(form_logout.submit())"><span class="glyphicon glyphicon-log-out" aria-hidden="true"></span> ログアウト</a></li>
							</ul>
						</li>
					
				</ul>
			</div>
		</div>
	</nav>

	<form method="POST" name="form_logout" action="/logout?continue=https%3A%2F%2Fatcoder.jp%2Fcontests%2Fabc001%2Ftasks%2Fabc001_3">
		<input type="hidden" name="csrf_token" value="J0H5Y5m/RVVl3DcGJmCaMbXU1XlSlNSccxiajEiH9pg=" />
	</form>
	<div id="main-container" class="container"
		 	style="padding-top:50px;">
		


<div class="row">
	<div id="contest-nav-tabs" class="col-sm-12 mb-2 cnvtb-fixed">
	<div>
		<small class="contest-duration">
			
				コンテスト時間:
				<a href='http://www.timeanddate.com/worldclock/fixedtime.html?iso=20131012T2100&p1=248' target='blank'><time class='fixtime fixtime-full'>2013-10-12 21:00:00+0900</time></a> ~ <a href='http://www.timeanddate.com/worldclock/fixedtime.html?iso=20131012T2300&p1=248' target='blank'><time class='fixtime fixtime-full'>2013-10-12 23:00:00+0900</time></a> 
				(120分)
			
		</small>
		<small class="back-to-home pull-right"><a href="/home">AtCoderホームへ戻る</a></small>
	</div>
	<ul class="nav nav-tabs">
		<li><a href="/contests/abc001"><span class="glyphicon glyphicon-home" aria-hidden="true"></span> トップ</a></li>
		
			<li class="active"><a href="/contests/abc001/tasks"><span class="glyphicon glyphicon-tasks" aria-hidden="true"></span> 問題</a></li>
		

		
			<li><a href="/contests/abc001/clarifications"><span class="glyphicon glyphicon-question-sign" aria-hidden="true"></span> 質問 <span id="clar-badge" class="badge"></span></a></li>
		

		
			<li><a href="/contests/abc001/submit?taskScreenName=abc001_3"><span class="glyphicon glyphicon-send" aria-hidden="true"></span> 提出</a></li>
		

		
			<li>
				<a class="dropdown-toggle" data-toggle="dropdown" href="#" role="button" aria-haspopup="true" aria-expanded="false"><span class="glyphicon glyphicon-list" aria-hidden="true"></span> 提出結果<span class="caret"></span></a>
				<ul class="dropdown-menu">
					<li><a href="/contests/abc001/submissions"><span class="glyphicon glyphicon-globe" aria-hidden="true"></span> すべての提出</a></li>
					
						<li><a href="/contests/abc001/submissions/me"><span class="glyphicon glyphicon-user" aria-hidden="true"></span> 自分の提出</a></li>
						<li class="divider"></li>
						<li><a href="/contests/abc001/score"><span class="glyphicon glyphicon-dashboard" aria-hidden="true"></span> 自分の得点状況</a></li>
					
				</ul>
			</li>
		

		
			
				
					<li><a href="/contests/abc001/standings"><span class="glyphicon glyphicon-sort-by-attributes-alt" aria-hidden="true"></span> 順位表</a></li>
				
			
				
					<li><a href="/contests/abc001/standings/virtual"><span class="glyphicon glyphicon-sort-by-attributes-alt" aria-hidden="true"></span> バーチャル順位表</a></li>
				
			
		

		
			<li><a href="/contests/abc001/custom_test"><span class="glyphicon glyphicon-wrench" aria-hidden="true"></span> コードテスト</a></li>
		

		
			<li><a href="/contests/abc001/editorial"><span class="glyphicon glyphicon-book" aria-hidden="true"></span> 解説</a></li>
		
		

		<li class="pull-right"><a id="fix-cnvtb" href="javascript:void(0)"><span class="glyphicon glyphicon-pushpin" aria-hidden="true"></span></a></li>
	</ul>
</div>
	<div class="col-sm-12">
		<span class="h2">
			C - 風力観測
			<a class="btn btn-default btn-sm" href="/contests/abc001/tasks/abc001_3/editorial">解説</a>
		</span>
		<span id="task-lang-btn" class="pull-right"><span data-lang="ja"><img src='//img.atcoder.jp/assets/top/img/flag-lang/ja.png'></span> / <span data-lang="en"><img src='//img.atcoder.jp/assets/top/img/flag-lang/en.png'></span></span>
		<script>
			$(function() {
				var ts = $('#task-statement span.lang');
				if (ts.children('span').size() <= 1) {
					$('#task-lang-btn').hide();
					ts.children('span').show();
					return;
				}
				var REMEMBER_LB = 5;
				var LS_KEY = 'task_lang';
				var taskLang = getLS(LS_KEY) || '';
				var changeTimes = 0;
				if (taskLang == 'ja' || taskLang == 'en') {
					changeTimes = REMEMBER_LB;
				} else {
					var changeTimes = parseInt(taskLang, 10);
					if (isNaN(changeTimes)) {
						changeTimes = 0;
						delLS(LS_KEY);
					}
					changeTimes++;
					taskLang = LANG;
				}
				ts.children('span.lang-' + taskLang).show();

				$('#task-lang-btn span').click(function() {
					var l = $(this).data('lang');
					ts.children('span').hide();
					ts.children('span.lang-' + l).show();
					if (changeTimes < REMEMBER_LB) setLS(LS_KEY, changeTimes);
					else setLS(LS_KEY, l);
				});
			});
		</script>
		<hr/>
		<p>
			実行時間制限: 2 sec / メモリ制限: 64 MB
			
		</p>

		<div id="task-statement">
			
			<div class="part">
<h4>注意</h4>
<p><b>この問題は古い問題です。過去問練習をする場合は、新しいAtCoder Beginner Contestから取り組むことをお勧めしています。</b></p>

<h3>問題文</h3>
<section>
ある風向風速計は、風向の角度と風程を <var>1</var> 分毎に自動で記録してくれます。<br />
<br />
風向の角度というのは真北を <var>0</var> 度とし、そこから時計回りに風の吹いてくる方向を定めたものです。気象観測等では全体を等しく <var>16</var> 分割した <var>16</var> 方位が用いられます。<var>16</var> 方位と角度は、以下の表のように対応します。<br />
<br />
<center>
<table> 
<caption> 風向と角度の関係</caption>
<tr>
<th>　方位　</th><th>　角度　</th>
<th>　方位　</th><th>　角度　</th>
</tr>
<tr>
<th>　N (北)　</th><th>　他のいずれにも当てはまらない　</th>
<th>　S (南)　</th><th>　<var>168.75</var> 度以上 <var>191.25</var> 度未満　</th>
</tr>
<tr>
<th>　NNE (北北東)　</th><th>　<var>11.25</var> 度以上 <var>33.75</var> 度未満　</th>
<th>　SSW (南南西)　</th><th>　<var>191.25</var> 度以上 <var>213.75</var> 度未満　</th>
</tr>
<tr>
<th>　NE (北東)　</th><th>　<var>33.75</var> 度以上 <var>56.25</var> 度未満　</th>
<th>　SW (南西)　</th><th>　<var>213.75</var> 度以上 <var>236.25</var> 度未満　</th>
</tr>
<tr>
<th>　ENE (東北東)　</th><th>　<var>56.25</var> 度以上 <var>78.75</var> 度未満　</th>
<th>　WSW (西南西)　</th><th>　<var>236.25</var> 度以上 <var>258.75</var> 度未満　</th>
</tr>
<tr>
<th>　E (東)　</th><th>　<var>78.75</var> 度以上 <var>101.25</var> 度未満　</th>
<th>　W (西)　</th><th>　<var>258.75</var> 度以上 <var>281.25</var> 度未満　</th>
</tr>
<tr>
<th>　ESE (東南東)　</th><th>　<var>101.25</var> 度以上 <var>123.75</var> 度未満　</th>
<th>　WNW (西北西)　</th><th>　<var>281.25</var> 度以上 <var>303.75</var> 度未満　</th>
</tr>
<tr>
<th>　SE (南東)　</th><th>　<var>123.75</var> 度以上 <var>146.25</var> 度未満　</th>
<th>　NW (北西)　</th><th>　<var>303.75</var> 度以上 <var>326.25</var> 度未満　</th>
</tr>
<tr>
<th>　SSE (南南東)　</th><th>　<var>146.25</var> 度以上 <var>168.75</var> 度未満　</th>
<th>　NNW (北北西)　</th><th>　<var>326.25</var> 度以上 <var>348.75</var> 度未満　</th>
</tr>
</table>
</center>
<br />
<br />
風程というのは、風向風速計の風車が、ある一定の時間に風によって回った量を距離によって表したものです。<br />
例えば、<var>1</var> 分間の風程が <var>300{\rm m}</var> というのは、<var>1</var> 分間に吹いた風によって風車が <var>300{\rm m}</var> 回ったという事です。この時、この <var>1</var> 分間の平均風速は <var>300{\rm m}</var> を <var>60</var> 秒で割って、<var>5{\rm m}/{\rm s}</var> と求められます。<br />
<br />
与えられたデータを、ラジオ等で流れる「気象通報」と同様の形式に直そうと思います。<br />
気象通報では、<var>16</var> 方位での風向と、風力 (ビューフォート風力階級) が伝えられます。<br />
<br />
風向は、先の表の <var>16</var> 方位です。
ただし、風力 <var>0</var> の場合、実際には「風弱く」と伝えられるため、風向は <var>16</var> 方位ではなく、特別な向きである<code>C</code>とします。<br />
<br />
風力は、風速を計算し、小数第 <var>2</var> 位を四捨五入して、以下の対応により風力に変換します。
<br />　
<center>
<table> 
<caption> 風力と風速の関係 　(ビューフォート風力階級)</caption>
<tr>
<th>風力　　</th><th>風速　　</th>
<th>風力　　</th><th>風速　　</th>
<th>風力　　</th><th>風速　　</th>
</tr>
<tr>
<th>風力<var>0</var>　　</th> <th><var>0.0{\rm m}/{\rm s}</var> 以上 <var>0.2{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>5</var>　　</th> <th><var>8.0{\rm m}/{\rm s}</var> 以上 <var>10.7{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>10</var>　　</th> <th><var>24.5{\rm m}/{\rm s}</var> 以上 <var>28.4{\rm m}/{\rm s}</var> 以下　　</th>
</tr>
<tr>
<th>風力<var>1</var>　　</th> <th><var>0.3{\rm m}/{\rm s}</var> 以上 <var>1.5{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>6</var>　　</th> <th><var>10.8{\rm m}/{\rm s}</var> 以上 <var>13.8{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>11</var>　　</th> <th><var>28.5{\rm m}/{\rm s}</var> 以上 <var>32.6{\rm m}/{\rm s}</var> 以下　　</th>
</tr>
<tr>
<th>風力<var>2</var>　　</th> <th><var>1.6{\rm m}/{\rm s}</var> 以上 <var>3.3{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>7</var>　　</th> <th><var>13.9{\rm m}/{\rm s}</var> 以上 <var>17.1{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>12</var>　　</th> <th><var>32.7{\rm m}/{\rm s}</var> 以上　　</th>
</tr>
<tr>
<th>風力<var>3</var>　　</th> <th><var>3.4{\rm m}/{\rm s}</var> 以上 <var>5.4{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>8</var>　　</th> <th><var>17.2{\rm m}/{\rm s}</var> 以上 <var>20.7{\rm m}/{\rm s}</var> 以下　　</th>
<th>　</th>
</tr>
<tr>
<th>風力<var>4</var>　　</th> <th><var>5.5{\rm m}/{\rm s}</var> 以上 <var>7.9{\rm m}/{\rm s}</var> 以下　　</th>
<th>風力<var>9</var>　　</th> <th><var>20.8{\rm m}/{\rm s}</var> 以上 <var>24.4{\rm m}/{\rm s}</var> 以下　　</th>
<th>　</th>
</tr>
</table>
</center>

<br />
風向 (角度) と <var>1</var> 分間の風程が入力されるとき、それを気象通報の形式に直して出力するプログラムを作成してください。
</section>
</div>

<hr />

<div class="io-style">
<div class="part">
<h3>入力</h3>
<section>
入力は以下の形式の <var>1</var> 行からなる。
<pre>
<var>Deg</var> <var>Dis</var>
</pre>
<ul>
<li><var>Deg</var>は風向を示し、本来の角度を <var>10</var> 倍した整数で与えられる (例えば、<var>90</var> 度なら <var>900</var>、<var>137.5</var> 度なら<var>1375</var>と与えられる) 。
</li>
<li>
<var>Dis</var>は <var>1</var> 分間の風程を示す整数である。単位はメートル <var>({\rm m})</var> である。
</li>
</ul>

</section>
</div>

<div class="part">
<h3>制約</h3>
<section>
<ul>
<li><var>0≦Deg<3,600</var></li>
<li><var>0≦Dis<12,000</var></li>
</ul>
</section>
</div>

<div class="part">
<h3>出力</h3>
<section>
出力は以下の形式の <var>1</var> 行とする。
<pre>
<var>Dir</var> <var>W</var>
</pre>
<ul>
<li>
<var>Dir</var>は風向を示し、<code>C</code>, <code>N</code>, <code>E</code>, <code>W</code>, <code>S</code> からなる <var>1</var>〜<var>3</var> 文字の文字列である。
</li>
<li>
<var>W</var>は風力を示し、<var>0</var> 以上 <var>12</var> 以下の整数である。
</li>
</ul>
また、出力の末尾には改行を入れること。
</section>
</div>
</div>
<hr />

<div class="part">
<h3>入力例 1</h3>
<section>
<pre class="prettyprint linenums">
2750 628
</pre>
<ul>
	<li>この場合、風向は <var>275</var> 度、風程は <var>1</var> 分あたり <var>628{\rm m}</var> である。</li>
</ul>
</section>
</div>

<div class="part">
<h3>出力例 1</h3>
<section>
<pre class="prettyprint linenums">
W 5
</pre>
<ul>
<li><var>275</var> 度は西向きなので、<code>W</code>と出力する。</li>
<li><var>1</var> 分で<var>628{\rm m}</var>ということは、<var>10.46{\rm m}/{\rm s}</var>なので、小数第 <var>2</var> 位を四捨五入して<var>10.5{\rm m}/{\rm s}</var>となり、風力 <var>5</var> に相当する。</li>
</ul>
</section>
</div>

<hr />

<div class="part">
<h3>入力例 2</h3>
<section>
<pre class="prettyprint linenums">
161 8
</pre>
</section>
</div>

<div class="part">
<h3>出力例 2</h3>
<section>
<pre class="prettyprint linenums">
C 0
</pre>
<ul>
	<li>風向は本来<code>NNE</code>だが、風力 <var>0</var> であるため<code>C</code>とする。</li>
</ul>
</section>

<hr />

<div class="part">
<h3>入力例 3</h3>
<section>
<pre class="prettyprint linenums">
3263 15
</pre>
</section>
</div>

<div class="part">
<h3>出力例 3</h3>
<section>
<pre class="prettyprint linenums">
NNW 1
</pre>
<ul>
	<li>浮動小数点数型での計算は、誤差が発生する恐れがあります。</li>
	<li>環境によって計算結果が変わることもありますので、誤差には十分注意しましょう。</li>
</ul>
</section>

</div>

<hr />

<div class="part">
<h3>入力例 4</h3>
<section>
<pre class="prettyprint linenums">
1462 1959
</pre>
</section>
</div>

<div class="part">
<h3>出力例 4</h3>
<section>
<pre class="prettyprint linenums">
SE 12
</pre>
<ul>
	<li>誤差には十分注意しましょう。</li>
</ul>
</section>

</div>

<hr />

<div class="part">
<h3>入力例 5</h3>
<section>
<pre class="prettyprint linenums">
1687 1029
</pre>
</section>
</div>

<div class="part">
<h3>出力例 5</h3>
<section>
<pre class="prettyprint linenums">
SSE 8
</pre>
<ul>
	<li>誤差には十分注意しましょう。</li>
</ul>
</section>

</div>

<hr />

<div class="part">
<h3>入力例 6</h3>
<section>
<pre class="prettyprint linenums">
2587 644
</pre>
</section>
</div>

<div class="part">
<h3>出力例 6</h3>
<section>
<pre class="prettyprint linenums">
WSW 5
</pre>
<ul>
	<li>誤差には十分注意しましょう。</li>
</ul>
</section>

</div>

<hr />

<div class="part">
<h3>入力例 7</h3>
<section>
<pre class="prettyprint linenums">
113 201
</pre>
</section>
</div>

<div class="part">
<h3>出力例 7</h3>
<section>
<pre class="prettyprint linenums">
NNE 3
</pre>
<ul>
	<li>誤差には十分注意しましょう。</li>
</ul>
</section>

</div>

<hr />

<div class="part">
<h3>入力例 8</h3>
<section>
<pre class="prettyprint linenums">
2048 16
</pre>
</section>
</div>

<div class="part">
<h3>出力例 8</h3>
<section>
<pre class="prettyprint linenums">
SSW 1
</pre>
</section>

</div>
</div>
		</div>

		

		
		<hr/>
		<form class="form-horizontal form-code-submit" action="/contests/abc001/submit" method="POST">
			<input type="hidden" name="data.TaskScreenName" value="abc001_3" />
			
			<div class="form-group ">
				<label class="control-label col-sm-2" for="select-lang">言語</label>
				<div id="select-lang" class="col-sm-5" data-name="data.LanguageId">
					<select class="form-control current" data-placeholder="-" name="data.LanguageId" required>
						<option></option>
						
							<option value="4001" data-mime="text/x-csrc">C (GCC 9.2.1)</option>
						
							<option value="4002" data-mime="text/x-csrc">C (Clang 10.0.0)</option>
						
							<option value="4003" data-mime="text/x-c&#43;&#43;src">C&#43;&#43; (GCC 9.2.1)</option>
						
							<option value="4004" data-mime="text/x-c&#43;&#43;src">C&#43;&#43; (Clang 10.0.0)</option>
						
							<option value="4005" data-mime="text/x-java">Java (OpenJDK 11.0.6)</option>
						
							<option value="4006" data-mime="text/x-python">Python (3.8.2)</option>
						
							<option value="4007" data-mime="text/x-sh">Bash (5.0.11)</option>
						
							<option value="4008" data-mime="text/x-bc">bc (1.07.1)</option>
						
							<option value="4009" data-mime="text/x-sh">Awk (GNU Awk 4.1.4)</option>
						
							<option value="4010" data-mime="text/x-csharp">C# (.NET Core 3.1.201)</option>
						
							<option value="4011" data-mime="text/x-csharp">C# (Mono-mcs 6.8.0.105)</option>
						
							<option value="4012" data-mime="text/x-csharp">C# (Mono-csc 3.5.0)</option>
						
							<option value="4013" data-mime="text/x-clojure">Clojure (1.10.1.536)</option>
						
							<option value="4014" data-mime="text/x-crystal">Crystal (0.33.0)</option>
						
							<option value="4015" data-mime="text/x-d">D (DMD 2.091.0)</option>
						
							<option value="4016" data-mime="text/x-d">D (GDC 9.2.1)</option>
						
							<option value="4017" data-mime="text/x-d">D (LDC 1.20.1)</option>
						
							<option value="4018" data-mime="application/dart">Dart (2.7.2)</option>
						
							<option value="4019" data-mime="text/x-dc">dc (1.4.1)</option>
						
							<option value="4020" data-mime="text/x-erlang">Erlang (22.3)</option>
						
							<option value="4021" data-mime="elixir">Elixir (1.10.2)</option>
						
							<option value="4022" data-mime="text/x-fsharp">F# (.NET Core 3.1.201)</option>
						
							<option value="4023" data-mime="text/x-fsharp">F# (Mono 10.2.3)</option>
						
							<option value="4024" data-mime="text/x-forth">Forth (gforth 0.7.3)</option>
						
							<option value="4025" data-mime="text/x-fortran">Fortran (GNU Fortran 9.2.1)</option>
						
							<option value="4026" data-mime="text/x-go">Go (1.14.1)</option>
						
							<option value="4027" data-mime="text/x-haskell">Haskell (GHC 8.8.3)</option>
						
							<option value="4028" data-mime="text/x-haxe">Haxe (4.0.3); js</option>
						
							<option value="4029" data-mime="text/x-haxe">Haxe (4.0.3); Java</option>
						
							<option value="4030" data-mime="text/javascript">JavaScript (Node.js 12.16.1)</option>
						
							<option value="4031" data-mime="text/x-julia">Julia (1.4.0)</option>
						
							<option value="4032" data-mime="text/x-kotlin">Kotlin (1.3.71)</option>
						
							<option value="4033" data-mime="text/x-lua">Lua (Lua 5.3.5)</option>
						
							<option value="4034" data-mime="text/x-lua">Lua (LuaJIT 2.1.0)</option>
						
							<option value="4035" data-mime="text/x-sh">Dash (0.5.8)</option>
						
							<option value="4036" data-mime="text/x-nim">Nim (1.0.6)</option>
						
							<option value="4037" data-mime="text/x-objectivec">Objective-C (Clang 10.0.0)</option>
						
							<option value="4038" data-mime="text/x-common-lisp">Common Lisp (SBCL 2.0.3)</option>
						
							<option value="4039" data-mime="text/x-ocaml">OCaml (4.10.0)</option>
						
							<option value="4040" data-mime="text/x-octave">Octave (5.2.0)</option>
						
							<option value="4041" data-mime="text/x-pascal">Pascal (FPC 3.0.4)</option>
						
							<option value="4042" data-mime="text/x-perl">Perl (5.26.1)</option>
						
							<option value="4043" data-mime="text/x-perl">Raku (Rakudo 2020.02.1)</option>
						
							<option value="4044" data-mime="text/x-php">PHP (7.4.4)</option>
						
							<option value="4045" data-mime="text/x-prolog">Prolog (SWI-Prolog 8.0.3)</option>
						
							<option value="4046" data-mime="text/x-python">PyPy2 (7.3.0)</option>
						
							<option value="4047" data-mime="text/x-python">PyPy3 (7.3.0)</option>
						
							<option value="4048" data-mime="text/x-racket">Racket (7.6)</option>
						
							<option value="4049" data-mime="text/x-ruby">Ruby (2.7.1)</option>
						
							<option value="4050" data-mime="text/x-rustsrc">Rust (1.42.0)</option>
						
							<option value="4051" data-mime="text/x-scala">Scala (2.13.1)</option>
						
							<option value="4052" data-mime="text/x-java">Java (OpenJDK 1.8.0)</option>
						
							<option value="4053" data-mime="text/x-scheme">Scheme (Gauche 0.9.9)</option>
						
							<option value="4054" data-mime="text/x-sml">Standard ML (MLton 20130715)</option>
						
							<option value="4055" data-mime="text/x-swift">Swift (5.2.1)</option>
						
							<option value="4056" data-mime="text/plain">Text (cat 8.28)</option>
						
							<option value="4057" data-mime="text/typescript">TypeScript (3.8)</option>
						
							<option value="4058" data-mime="text/x-vb">Visual Basic (.NET Core 3.1.101)</option>
						
							<option value="4059" data-mime="text/x-sh">Zsh (5.4.2)</option>
						
							<option value="4060" data-mime="text/x-cobol">COBOL - Fixed (OpenCOBOL 1.1.0)</option>
						
							<option value="4061" data-mime="text/x-cobol">COBOL - Free (OpenCOBOL 1.1.0)</option>
						
							<option value="4062" data-mime="text/x-brainfuck">Brainfuck (bf 20041219)</option>
						
							<option value="4063" data-mime="text/x-ada">Ada2012 (GNAT 9.2.1)</option>
						
							<option value="4064" data-mime="text/x-unlambda">Unlambda (2.0.0)</option>
						
							<option value="4065" data-mime="text/x-python">Cython (0.29.16)</option>
						
							<option value="4066" data-mime="text/x-sh">Sed (4.4)</option>
						
							<option value="4067" data-mime="text/x-vim">Vim (8.2.0460)</option>
						
					</select>
					<span class="error"></span>
				</div>
			</div>
			<script>var currentLang = getLS('defaultLang');</script>
			
			
<div class="form-group">
	<label class="control-label col-sm-2" for="sourceCode">ソースコード</label>
	<div class="col-sm-7" id="sourceCode">
		<div class="div-editor">
			<textarea class="form-control editor" name="sourceCode"></textarea>
		</div>
		<textarea class="form-control plain-textarea" style="display:none;"></textarea>
		<p>
			<span class="gray">※ 512 KiB まで</span><br>
			<span class="gray">※ ソースコードは「Main.<i>拡張子</i>」で保存されます</span>
		</p>
	</div>
	<div class="col-sm-3 editor-buttons">
		<p><button id="btn-open-file" type="button" class="btn btn-default btn-sm">
			<span class="glyphicon glyphicon-folder-open" aria-hidden="true"></span> &nbsp; ファイルを開く
		</button></p>
		<p><button type="button" class="btn btn-default btn-sm btn-toggle-editor" data-toggle="button" aria-pressed="false" autocomplete="off">
			エディタ切り替え
		</button></p>
		<p><button type="button" class="btn btn-default btn-sm btn-auto-height" data-toggle="button" aria-pressed="false" autocomplete="off">
			高さ自動調節
		</button></p>
	</div>
	<input id="input-open-file" type="file" style="display:none;">
</div>

			<input type="hidden" name="csrf_token" value="J0H5Y5m/RVVl3DcGJmCaMbXU1XlSlNSccxiajEiH9pg=" />
			<div class="form-group">
				<label class="control-label col-sm-2" for="submit"></label>
				<div class="col-sm-5">
					<button type="submit" class="btn btn-primary" id="submit">提出</button>
				</div>
			</div>
		</form>
		
	</div>
</div>




		
			<hr>
			
			
			
<div class="a2a_kit a2a_kit_size_20 a2a_default_style pull-right" data-a2a-url="https://atcoder.jp/contests/abc001/tasks/abc001_3?lang=ja" data-a2a-title="C - 風力観測">
	<a class="a2a_button_facebook"></a>
	<a class="a2a_button_twitter"></a>
	
		<a class="a2a_button_hatena"></a>
	
	<a class="a2a_dd" href="https://www.addtoany.com/share"></a>
</div>

		
		<script async src="//static.addtoany.com/menu/page.js"></script>
		
	</div> 
	<hr>
</div> 

	<div class="container" style="margin-bottom: 80px;">
			<footer class="footer">
			
				<ul>
					<li><a href="/contests/abc001/rules">ルール</a></li>
					<li><a href="/contests/abc001/glossary">用語集</a></li>
					
				</ul>
			
			<ul>
				<li><a href="/tos">利用規約</a></li>
				<li><a href="/privacy">プライバシーポリシー</a></li>
				<li><a href="/personal">個人情報保護方針</a></li>
				<li><a href="/company">企業情報</a></li>
				<li><a href="/faq">よくある質問</a></li>
				<li><a href="/contact">お問い合わせ</a></li>
				<li><a href="/documents/request">資料請求</a></li>
			</ul>
			<div class="text-center">
					<small id="copyright">Copyright Since 2012 &copy;<a href="http://atcoder.co.jp">AtCoder Inc.</a> All rights reserved.</small>
			</div>
			</footer>
	</div>
	<p id="fixed-server-timer" class="contest-timer"></p>
	<div id="scroll-page-top" style="display:none;"><span class="glyphicon glyphicon-arrow-up" aria-hidden="true"></span> ページトップ</div>

</body>
</html>


`

	extractor := NewFullTextExtractor()
	ja, en, err := extractor.Extract(strings.NewReader(source))
	if err != nil {
		t.Errorf("failed to extract statement: %s", err.Error())
	}

	t.Errorf("ja: %v, len(ja): %d, en: %v, len(en): %d", ja, len(ja), en, len(en))
}
