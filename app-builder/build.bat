SET APP=led-light-strip
docker run --rm -v "%cd%":/usr/src/%APP% --platform linux/arm/v6 -w /usr/src/%APP% ws2811-builder:latest go build -o "%APP%-armv6" -v
