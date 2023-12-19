run:
  bash ./scripts/tailwind.sh -i input.scss -o ./static/styles.css
  ~/go/bin/templ generate
  go run .
