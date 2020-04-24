#!/bin/sh -e
# Generate assets

DEBUG=$1

# Icons from svg
for name in img/*.svg; do
    echo "Generate ${name}"
    newname=$(echo "${name}" | sed 's/.svg//')
    convert -background none -density 553 -resize "64x64" "${name}" "gui/default/assets/${newname}.generated.png"
done

get() {
    url=$1
    sha1=$2
    out_file=$3
    echo "Downloading ${url}"
    echo "${sha1} -" > /tmp/sha1-sum.txt
    mkdir -p "$(dirname "${out_file}")"
    wget -qO- "${url}" | tee "${out_file}" | sha1sum -c /tmp/sha1-sum.txt
    rm -f /tmp/sha1-sum.txt
}

get_tar() {
    url=$1
    sha1=$2
    out_dir=$3
    shift; shift; shift
    get "${url}" "${sha1}" /tmp/data.tar.gz
    rm -rf "${out_dir}"; mkdir -p "${out_dir}"
    tar xf /tmp/data.tar.gz -C "${out_dir}" --wildcards $@
    rm -f /tmp/data.tar.gz
}

# Download gui vendor libraries
get_tar https://github.com/ForkAwesome/Fork-Awesome/archive/1.1.7.tar.gz \
    e5051a8c9b00ae1c6e0cf8958150f6cce952badf gui/default/vendor/fork-awesome \
    --strip-components=1 "Fork-Awesome-1.1.7/css/*.min.css" Fork-Awesome-1.1.7/fonts

get https://ajax.googleapis.com/ajax/libs/angularjs/1.7.9/angular.min.js \
    73b623b7d29122a34e73a061491f708b3b7f9f83 gui/default/vendor/angular/angular.min.js
get https://cdnjs.cloudflare.com/ajax/libs/ngStorage/0.3.11/ngStorage.min.js \
    95eece4152f3eda1aa5f626897f5dd8c790bfa2e gui/default/vendor/angular/ngStorage.min.js
get https://code.angularjs.org/1.7.9/angular-route.min.js \
    a54d31d32e3135ac6c7e555796bd213ed1ebd303 gui/default/vendor/angular/angular-route.min.js

get https://code.jquery.com/jquery-3.5.0.min.js \
    1d6ae46f2ffa213dede37a521b011ec1cd8d1ad3 gui/default/vendor/jquery/jquery.min.js

get https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css \
    b0972fdcce82fd583d4c2ccc3f2e3df7404a19d0 gui/default/vendor/bootstrap/css/bootstrap.min.css
get https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap-theme.min.css \
    beee0e080ea6dcc8c0661b66c1baa08e45f4ecb6 gui/default/vendor/bootstrap/css/bootstrap-theme.min.css
get https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js \
    b16fd8226bd6bfb08e568f1b1d0a21d60247cefb gui/default/vendor/bootstrap/js/bootstrap.min.js
