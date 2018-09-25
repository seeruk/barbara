package barbara

// BuildStyleSheet ...
func BuildStyleSheet() string {
	// TODO(elliot): Build from config, using text/template?
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
