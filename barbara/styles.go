package barbara

// BuildStyleSheet returns the application's stylesheet.
// TODO(elliot): When we have an Application type, this should be a method of that type.
// TODO(elliot): Build from config, using text/template?
func BuildStyleSheet() string {
	return `
		QMainWindow {
			background: #1a1a1a;
			margin: 0;
			padding: 0px;
		}

		QLabel {
			color: #e5e5e5;
			font-family: "Fira Sans";
			font-size: 13px;
			padding: 0 0 0 7px;
			text-align: center;
		}

		.barbara-button {
			background-color: #1a1a1a;
			color: #e5e5e5;
			font-family: "Fira Sans";
			font-size: 13px;
			padding: 7px;
		}

		.barbara-button:flat {
			border: 1px solid #333;
			border-radius: 3px;
		}
	`
}
