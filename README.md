# Hermes

[![Build Status](https://github.com/chauhannishith/hermes/actions/workflows/main.yml/badge.svg)](https://github.com/chauhannishith/hermes/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chauhannishith/hermes/v2)](https://goreportcard.com/report/github.com/chauhannishith/hermes/v2)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/chauhannishith/hermes/v2.svg)](https://pkg.go.dev/github.com/chauhannishith/hermes/v2)

> **Fork:** Maintained fork of [matcornic/hermes](https://github.com/matcornic/hermes) with optional greeting, signature, and unsubscribe controls. Install with `go get github.com/chauhannishith/hermes/v2@v2.1.1`.

Hermes is the Go port of the great [mailgen](https://github.com/eladnava/mailgen) engine for Node.js. Check their work, it's awesome!
It's a package that generates clean, responsive HTML e-mails for sending transactional e-mails (welcome e-mails, reset password e-mails, receipt e-mails and so on), and associated plain text fallback.

# Demo

<img src="screens/default/welcome.png" height="400" /> <img src="screens/default/reset.png" height="400" /> <img src="screens/default/receipt.png" height="400" />

# Usage

First install the package:

```
go get github.com/chauhannishith/hermes/v2@v2.1.1
```

Import path:

```go
import "github.com/chauhannishith/hermes/v2"
```

## Migrate from upstream `matcornic/hermes`

This fork is a **v2 module** (`github.com/chauhannishith/hermes/v2`). If you were on the original package:

```bash
go get github.com/chauhannishith/hermes/v2@v2.1.1
```

Replace imports:

- `github.com/matcornic/hermes` → `github.com/chauhannishith/hermes/v2`
- `github.com/matcornic/hermes/v2` → `github.com/chauhannishith/hermes/v2`

## Use Hermes

Then, start using the package by importing and configuring it:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes{
    // Optional Theme
    // Theme: new(Default) 
    Product: hermes.Product{
        // Appears in header & footer of e-mails
        Name: "Hermes",
        Link: "https://example-hermes.com/",
        // Optional product logo
        Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
    },
}
```

Next, generate an e-mail using the following code:

```go
email := hermes.Email{
    Body: hermes.Body{
        Name: "Jon Snow",
        Intros: []string{
            "Welcome to Hermes! We're very excited to have you on board.",
        },
        Actions: []hermes.Action{
            {
                Instructions: "To get started with Hermes, please click here:",
                Button: hermes.Button{
                    Color: "#22BC66", // Optional action button color
                    Text:  "Confirm your account",
                    Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
                },
            },
        },
        Outros: []string{
            "Need help, or have questions? Just reply to this email, we'd love to help.",
        },
    },
}

// Generate an HTML email with the provided contents (for modern clients)
emailBody, err := h.GenerateHTML(email)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}

// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
emailText, err := h.GeneratePlainText(email)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}

// Optionally, preview the generated HTML e-mail by writing it to a local file
err = os.WriteFile("preview.html", []byte(emailBody), 0644)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}
```

This code would output the following HTML template:

<img src="screens/demo.png" height="400" />

And the following plain text:

```

------------
Hi Jon Snow,
------------

Welcome to Hermes! We're very excited to have you on board.

To get started with Hermes, please click here: https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010

Need help, or have questions? Just reply to this email, we'd love to help.

Yours truly,
Hermes - https://example-hermes.com/

Copyright © 2025 Hermes. All rights reserved.
```

> Theme templates are embedded in your application binary. To use your own layout or load templates from files, implement a [custom theme](#custom-themes).

## More Examples

* [Welcome with button](examples/welcome.go)
* [Welcome with invite code](examples/invite_code.go)
* [Receipt](examples/receipt.go)
* [Password Reset](examples/reset.go)
* [Maintenance](examples/maintenance.go)
* [Sample custom CSS overrides](examples/custom.css) — see [Custom CSS overrides](#custom-css-overrides) for usage and testing

To run the examples, go to `examples` folder, then run `go run .`. HTML and plaintext examples are created in `examples/default/`.

Optionaly you can set the following variables to send automatically the emails to one your mailbox. Nice for testing template in real email clients.

* `HERMES_SEND_EMAILS=true`
* `HERMES_SMTP_SERVER=<smtp_server>` : for Gmail it's `smtp.gmail.com`
* `HERMES_SMTP_PORT=<smtp_port>` : for Gmail it's `465`
* `HERMES_SENDER_EMAIL=<your_sender_email>`
* `HERMES_SENDER_IDENTITY=<the sender name>`
* `HERMES_SMTP_USER=<smtp user>` : usually the same than `HERMES_SENDER_EMAIL`
* `HERMES_TO=<recipients emails>`: split by commas like `myadress@test.com,somethingelse@gmail.com`

The program will ask for your SMTP password. If needed, you can set it with `HERMES_SMTP_PASSWORD` variable (but be careful where you put this information !)

## Plaintext E-mails

To generate a [plaintext version of the e-mail](https://litmus.com/blog/best-practices-for-plain-text-emails-a-look-at-why-theyre-important), simply call `GeneratePlainText` function:

```go
// Generate plaintext email using hermes
emailText, err := h.GeneratePlainText(email)
if err != nil {
    panic(err) // Tip: Handle error with something else than a panic ;)
}
```

## Supported Themes

The following open-source themes are bundled with this package:

* `default` by [Postmark Transactional Email Templates](https://github.com/wildbit/postmark-templates)

<img src="screens/default/welcome.png" height="200" /> <img src="screens/default/reset.png" height="200" /> <img src="screens/default/receipt.png" height="200" />

### Custom themes

Only the `default` theme ships with this package. To change the **layout** (not just colors or fonts), implement `hermes.Theme` and pass it when creating Hermes:

```go
h := hermes.Hermes{
    Theme: new(MyCustomTheme), // see default.go for the interface
    Product: hermes.Product{
        Name: "Acme",
        Link: "https://example.com/",
    },
}
```

Your theme must provide HTML and plain-text templates plus base styles. Templates receive `Hermes`, `Email`, and `StylesCSS` — use the injection snippets so `Body` fields (intros, actions, table, etc.) render correctly.

**Full guide:** [Using a Custom Theme](CONTRIBUTING.md#using-a-custom-theme) in [CONTRIBUTING.md](CONTRIBUTING.md) (interface methods, template snippets, built-in theme checklist).

For **styling only** on the default layout, use [`CustomCSS`](#custom-css-overrides) instead — no custom theme required.

## RTL Support

To change the default text direction (left-to-right), simply override it as follows:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes {
    // Custom text direction
    TextDirection: hermes.TDRightToLeft,
}
```

## Language Customizations

Set a default greeting for all emails on `Hermes`, or override it per email on `Body`. The default greeting is `Hi`.

```go
h := hermes.Hermes{
    DefaultGreeting: "Dear",
    Product: hermes.Product{
        Name: "Acme",
        Link: "https://example.com/",
    },
}
```

Per-email overrides:

```go
email := hermes.Email{
    Body: hermes.Body{
        Greeting:  "Hello",
        Signature: "Sincerely",
    },
}
```

To omit the greeting from the title line, set `DisableGreeting` to `true`. The recipient name still appears unless `Title` is set:

```go
email := hermes.Email{
    Body: hermes.Body{
        DisableGreeting: true,
    },
}
```

To omit the signature line entirely, set `DisableSignature` to `true`. The product name still appears in the closing block:

```go
email := hermes.Email{
    Body: hermes.Body{
        DisableSignature: true,
    },
}
```

To add an optional unsubscribe link in the email footer, set `UnsubscribeLink` on `Body`. Use `UnsubscribeText` to customize the link label (defaults to `Unsubscribe`):

```go
email := hermes.Email{
    Body: hermes.Body{
        UnsubscribeLink: "https://example.com/unsubscribe?token=abc",
        UnsubscribeText: "Manage email preferences", // optional
    },
}
```

To use a custom title string rather than a greeting/name introduction, provide it instead of `Name`:

```go
email := hermes.Email{
    Body: hermes.Body{
        // Title will override `Name`
        Title: "Welcome to Hermes",
    },
}
```

To customize the `Copyright`, override it when initializing `Hermes` within your `Product` as follows. If omitted, the default copyright uses the current year automatically:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes{
    // Optional Theme
    // Theme: new(Default)
    Product: hermes.Product{
        // Appears in header & footer of e-mails
        Name: "Hermes",
        Link: "https://example-hermes.com/",
        // Custom copyright notice (year is auto-filled in the default when Copyright is omitted)
        Copyright: "Copyright © 2025 Dharma Initiative. All rights reserved."
    },
}
```

To use a custom fallback text at the end of the email, change the `TroubleText` field of the `hermes.Product` struct. The default value is `If you’re having trouble with the button '{ACTION}', copy and paste the URL below into your web browser.`. The `{ACTION}` placeholder will be replaced with the corresponding text of the supplied action button:

```go
// Configure hermes by setting a theme and your product info
h := hermes.Hermes{
    // Optional Theme
    // Theme: new(Default)
    Product: hermes.Product{
        // Custom trouble text
        TroubleText: "If the {ACTION}-button is not working for you, just copy and paste the URL below into your web browser."
    },
}
```

Since `v2.1.0`, Hermes is automatically inlining all CSS to improve compatibility with email clients, thanks to [Premailer](https://github.com/vanng822/go-premailer/premailer).
You can disable this feature by setting `DisableCSSInlining` of `Hermes` struct to `true`.

```go
h := hermes.Hermes{
    ...
    DisableCSSInlining: true,
}
```

### Custom CSS overrides

Hermes keeps a **fixed email layout** (header, body sections, footer). You customize **appearance** with CSS; **content** is set only through the existing API — `Hermes.Product`, `Body.Name`, `Body.Intros`, `Body.Actions`, and the other fields documented below. There is no per-email HTML template override.

Set `CustomCSS` on `Body` to load a stylesheet from a **file path**, **URL**, or as **inline CSS**. Only selectors you define are merged on top of the default theme — the structure and which values appear stay the same.

```go
email := hermes.Email{
    Body: hermes.Body{
        Name:      "Jon Snow",
        Intros:    []string{"Welcome!"},
        CustomCSS: "brand.css", // restyle .button, body, h1, etc.
    },
}
```

Use the class names from the default theme (e.g. `.button`, `.email-body`, `.content-cell`) in your CSS file. See [`examples/custom.css`](examples/custom.css).

For a completely different layout, see [Custom themes](#custom-themes) or [CONTRIBUTING.md — Using a Custom Theme](CONTRIBUTING.md#using-a-custom-theme).

**Test that overrides work**

```bash
# Unit/integration tests (merge logic + example file)
go test ./... -run CustomCSS

# Generate preview HTML under examples/custom/ (from repo root)
go test ./examples/ -run TestGenerateCustomCSSExamples
```

Open `examples/custom/custom.welcome.html` in a browser and compare with `examples/default/default.welcome.html`.

`go run .` in `examples/` still generates only the default theme — custom previews are produced by the test above.

## Elements

Hermes supports injecting custom elements such as dictionaries, tables and action buttons into e-mails.

### Action

To inject an action button in to the e-mail, supply the `Actions` object as follows:

```go
email := hermes.Email{
    Body: hermes.Body{
        Actions: []hermes.Action{
            {
                Instructions: "To get started with Hermes, please click here:",
                Button: hermes.Button{
                    Color: "#22BC66", // Optional action button color
                    Text:  "Confirm your account",
                    Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
                },
            },
        },
    },
}
```

Alternatively, instead of having a button, an action can be an invite code as follows:

```go
email := hermes.Email{
    Body: hermes.Body{
        Actions: []hermes.Action{
            {
                Instructions: "To get started with Hermes, please use the invite code:",
                InviteCode: "123456",
            },
        },
    },
}
```

To inject multiple action buttons in to the e-mail, supply another struct in Actions slice `Action`.

### Table

To inject a table into the e-mail, supply the `Table` object as follows:

```go
email := hermes.Email{
    Body: hermes.Body{
        Table: hermes.Table{
            Data: [][]hermes.Entry{
                // List of rows
                {   
                    // Key is the column name, Value is the cell value
                    // First object defines what columns will be displayed
                    {Key: "Item", Value: "Golang"},
                    {Key: "Description", Value: "Open source programming language that makes it easy to build simple, reliable, and efficient software"},
                    {Key: "Price", Value: "$10.99"},
                },
                {
                    {Key: "Item", Value: "Hermes"},
                    {Key: "Description", Value: "Programmatically create beautiful e-mails using Golang."},
                    {Key: "Price", Value: "$1.99"},
                },
            },
            Columns: hermes.Columns{
                // Custom style for each rows
                CustomWidth: map[string]string{
                    "Item":  "20%",
                    "Price": "15%",
                },
                CustomAlignment: map[string]string{
                    "Price": "right",
                },
            },
        },
    },
}
```

### Dictionary

To inject key-value pairs of data into the e-mail, supply the `Dictionary` object as follows:

```go
email := hermes.Email{
    Body: hermes.Body{
        Dictionary: []hermes.Entry{
            {Key: "Date", Value: "20 November 1887"},
            {Key: "Address", Value: "221B Baker Street, London"},
        },
    },
}
```

### Free Markdown

If you need more flexibility in the content of your generated e-mail, while keeping the same format than any other e-mail, use Markdown content. Supply the `FreeMarkdown` object as follows:

```go
email := hermes.Email{
		Body: hermes.Body{
			FreeMarkdown: `
> _Hermes_ service will shutdown the **1st August 2025** for maintenance operations. 

Services will be unavailable based on the following schedule:

| Services | Downtime |
| :------:| :-----------: |
| Service A | 2AM to 3AM |
| Service B | 4AM to 5AM |
| Service C | 5AM to 6AM |

---

Feel free to contact us for any question regarding this matter at [support@hermes-example.com](mailto:support@hermes-example.com) or in our [Gitter](https://gitter.im/)

`,
		},
	}
}
```

This code would output the following HTML template:

<img src="screens/free-markdown.png" height="400" />

And the following plaintext:

```
------------
Hi Jon Snow,
------------

> 
> 
> 
> Hermes service will shutdown the *1st August 2025* for maintenance
> operations.
> 
> 

Services will be unavailable based on the following schedule:

+-----------+------------+
| SERVICES  |  DOWNTIME  |
+-----------+------------+
| Service A | 2AM to 3AM |
| Service B | 4AM to 5AM |
| Service C | 5AM to 6AM |
+-----------+------------+

Feel free to contact us for any question regarding this matter at support@hermes-example.com ( support@hermes-example.com ) or in our Gitter ( https://gitter.im/ )

Yours truly,
Hermes - https://example-hermes.com/

Copyright © 2025 Hermes. All rights reserved.
```

Be aware that this content will replace existing tables, dictionary and actions. Only intros, outros, header and footer will be kept.

This is helpful when your application needs sending e-mails, wrote on-the-fly by adminstrators.

> Markdown is rendered with [Blackfriday](https://github.com/russross/blackfriday), so every thing Blackfriday can do, Hermes can do it as well.

## Troubleshooting

1. After sending multiple e-mails to the same Gmail / Inbox address, they become grouped and truncated since they contain similar text, breaking the responsive e-mail layout.

> Simply sending the `X-Entity-Ref-ID` header with your e-mails will prevent grouping / truncation.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). To add your own email layout, start with [Using a Custom Theme](CONTRIBUTING.md#using-a-custom-theme).

## Continuous integration

This fork runs [GitHub Actions](.github/workflows/) on every push and pull request (build, test, and lint). If workflows do not run after forking, enable them under **Settings → Actions → General** on your GitHub fork.

## License

Apache 2.0
