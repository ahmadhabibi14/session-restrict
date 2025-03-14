package mailer

func HtmlOtpNewSessionLoggedIn(
	title, resetPasswordLink, userName, timestamp,
	device, ip, accessToken string,
) string {
	return `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<meta name="color-scheme" content="light dark">
    <meta name="supported-color-schemes" content="light dark">
    <title>` + title + `</title>
    <style>
      *,
      *::before,
      *::after {
        box-sizing: border-box;
      }
      :root {
        -moz-tab-size: 4;
        -o-tab-size: 4;
        tab-size: 4;
        --blue-primary: #1187c9;
      }
      *::selection {
        background: var(--blue-primary);
        color: #fff;
      }
      html {
        max-width: 100%;
        margin: auto;
        line-height-step: 1.15;
        -webkit-text-size-adjust: 100%;
      }
      body {
        font-family: Arial, Roboto Helvetica, sans-serif;
        margin: 0;
        color: #121212;
      }
      main {
        padding: 20px;
      }
      picture img {
        width: 50%;
        height: auto;
        padding-top: 10px;
      }
      a {
        color: var(--blue-primary);
        text-decoration: none;
      }
      a:hover {
        text-decoration: underline;
      }
      .block {
        height: 5px;
        background-color: var(--blue-primary);
      }
      @media (prefers-color-scheme: dark) {
        body {
          background-color: #121212;
          color: #fff;
        }
      }
    </style>
  </head>
  <body
    style="
      font-family: Arial, Roboto Helvetica, sans-serif;
			font-size: 16px;
      color: #121212;
      background-color: #fff;
      margin: 0;
    "
  >
		<article style="max-width: 600px; margin: auto">
			<div style="height: 5px; background-color: #1187c9"></div>
			<main style="padding: 20px 0 0 0">
				<div style="padding: 0 10px; margin-top: 20px">
					<b>Halo ` + userName + `!!</b>
          <p>Kami mendeteksi login baru ke akun Anda:
          <p>
						📅 Waktu : ` + timestamp + `<br>
						🖥 Perangkat : ` + device + `<br>
						🌐 Alamat IP : ` + ip + `<br>
						🔑 Akses Token : ` + accessToken + `
					</p>
					<p>Jika ini memang Anda, tidak ada yang perlu dilakukan. Namun, jika ini bukan Anda, segera amankan akun Anda dengan:</p>
					<center>
						<a
							class="` + resetPasswordLink + `"
							style="border-radius: 10px; padding: 8px 15px; background-color: #1187c9; color: #ffffff; font-weight: 600"
						>
							Reset Password
						</a>
					</center>
				</div>
			</main>
			<div style="height: 5px; background-color: #1187c9"></div>
		</article>
  </body>
</html>
`
}
