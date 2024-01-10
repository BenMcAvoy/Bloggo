run:
  bun install
  bash ./scripts/tailwind.sh -i input.scss -o ./static/styles.css
  templ generate
  go run .
