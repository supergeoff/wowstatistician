<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		<title>Wow Statistician</title>
		<link
			rel="stylesheet"
			href="https://unpkg.com/bulmaswatch/lumen/bulmaswatch.min.css"
		/>
		<script
			defer
			src="https://use.fontawesome.com/releases/v5.3.1/js/all.js"
		></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.9.3/Chart.js"></script>
		<script src="/static/script/wasm_client/wasm_exec.js"></script>
		<script>
			async function init() {
				const go = new Go();
				let prog = await WebAssembly.instantiateStreaming(
					fetch("/static/script/wasm_client/wasm_client"),
					go.importObject
				);
				go.run(prog.instance);
			}
			init();
		</script>
	</head>
	<body>
		<section class="section">
			<div class="container">
				<div class="tile is-ancestor">
					<div class="tile is-parent is-vertical">
						<div class="tile is-child box">
							<div class="field is-horizontal">
								<div class="field-label">
									<label class="label">Show stats for:</label>
								</div>
								<div class="field-body">
									<div class="field">
										<div class="control is-expanded">
											<div class="select">
												<select id="source">
													<option>Mythic</option>
													<option>Raid</option>
													<option>Arena</option>
													<option>Rbg</option>
												</select>
											</div>
										</div>
										<p class="help">
											Show class and spec distribution
											based on Blizzard's leatherboard API
										</p>
									</div>
									<div class="field">
										<div class="control is-expanded">
											<label class="radio">
												<input
													type="radio"
													name="merge"
													value="spec"
												/>
												Spec
											</label>
											<label class="radio">
												<input
													type="radio"
													name="merge"
													value="class"
													checked
												/>
												Class
											</label>
										</div>
									</div>
								</div>
							</div>
							<div><canvas id="stats" height="400"></canvas></div>
						</div>
					</div>
				</div>
			</div>
			<footer class="footer">
				<div class="content has-text-centered">
					<p>
						<strong>Wow Statistician</strong>. The source code can
						be found on
						<a href="https://github.com/supergeoff/wowstatistician"
							>Github</a
						>. The data have been last updated on
						<span class="has-text-link" id="syncdate"></span>.
					</p>
				</div>
				<div class="content">
					<p
						class="is-size-7 has-text-grey-dark is-italic has-text-centered"
					>
						Tech Stack
					</p>
					<div class="tags is-centered">
						<a class="tag is-info" href="https://golang.org/">Go</a>
						<a class="tag is-info" href="https://webassembly.org/"
							>WASM</a
						>
						<a
							class="tag is-danger"
							href="https://github.com/dgraph-io/badger/"
							>Badger DB</a
						>
						<a class="tag is-danger" href="https://beego.me/"
							>Beego</a
						>
						<a class="tag is-danger" href="https://bulma.io/"
							>Bulma</a
						>
						<a class="tag is-danger" href="https://www.chartjs.org/"
							>ChartJS</a
						>
					</div>
				</div>
			</footer>
		</section>
	</body>
	<!-- <script src="/static/script/js_client/js_client.js"></script> -->
	<!-- <script src="/static/script/example/client.js"></script> -->
</html>
