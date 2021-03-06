<div class="fullPage">
    <div class="contentWrap">
        <img src="{{mountpathed "static/fantom_logo_white.svg"}}" alt="Fantom" />
        <form class="mainForm" action="{{mountpathed "register"}}" method="post">
            {{with .errors}}{{with (index . "")}}{{range .}}<span>{{.}}</span><br />{{end}}{{end}}{{end -}}
            <label for="name">Name:</label>
            <input name="name" type="text" value="{{with .preserve}}{{with .name}}{{.}}{{end}}{{end}}" placeholder="Name" /><br />
            {{with .errors}}{{range .name}}<span>{{.}}</span><br />{{end}}{{end -}}
            <label for="email">E-mail:</label>
            <input name="email" type="text" value="{{with .preserve}}{{with .email}}{{.}}{{end}}{{end}}" placeholder="E-mail" /><br />
            {{with .errors}}{{range .email}}<span>{{.}}</span><br />{{end}}{{end -}}
            <label for="password">Password:</label>
            <input name="password" type="password" placeholder="Password" /><br />
            {{with .errors}}{{range .password}}<span>{{.}}</span><br />{{end}}{{end -}}
            <label for="confirm_password">Confirm Password:</label>
            <input name="confirm_password" type="password" placeholder="Confirm Password" /><br />
            {{with .errors}}{{range .confirm_password}}<span>{{.}}</span><br />{{end}}{{end -}}
            <input class="loginBtn" type="submit" value="Register"><br />
            <a href="{{mountpathed "login"}}{{with .challenge}}?login_challenge={{.}}{{end}}">Cancel</a>
            {{with .csrf_token}}<input type="hidden" name="csrf_token" value="{{.}}" />{{end}}
            {{with .challenge}}<input type="hidden" name="challenge" value="{{.}}" />{{end}}
        </form>
    </div>
</div>
