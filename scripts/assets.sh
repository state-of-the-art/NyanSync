#!/bin/sh -e
# Generate assets

# Icons from svg
for name in img/*.svg; do
    echo "Generate ${name}"
    newname=$(echo "${name}" | sed 's/.svg//')
    convert -background none -density 553 -resize "64x64" "${name}" "gui/default/assets/${newname}.generated.png"
done

get_tar() {
    url=$1
    sha1=$2
    out_dir=$3
    shift; shift; shift
    echo "Downloading ${url}"
    echo "${sha1} -" > /tmp/sha1-sum.txt
    wget -qO- "${url}" | tee /tmp/data.tar.gz | sha1sum -c /tmp/sha1-sum.txt
    rm -rf "${out_dir}"; mkdir -p "${out_dir}"
    tar xf /tmp/data.tar.gz -C "${out_dir}" --wildcards $@
}

get() {
    url=$1
    sha1=$2
    out_file=$3
    echo "Downloading ${url}"
    echo "${sha1} -" > /tmp/sha1-sum.txt
    mkdir -p "$(dirname "${out_file}")"
    wget -qO- "${url}" | tee "${out_file}" | sha1sum -c /tmp/sha1-sum.txt
}

# Download gui vendor libraries
get_tar https://github.com/ForkAwesome/Fork-Awesome/archive/1.1.7.tar.gz \
    e5051a8c9b00ae1c6e0cf8958150f6cce952badf gui/default/vendor/fork-awesome \
    --strip-components=1 "Fork-Awesome-1.1.7/css/*.min.css" Fork-Awesome-1.1.7/fonts

get https://ajax.googleapis.com/ajax/libs/angularjs/1.7.9/angular.min.js \
    73b623b7d29122a34e73a061491f708b3b7f9f83 gui/default/vendor/angular/angular.min.js
get https://ajax.googleapis.com/ajax/libs/angularjs/1.7.9/angular-sanitize.min.js \
    a791888305092420a1aa859197ade2698766020d gui/default/vendor/angular/angular-sanitize.min.js
get https://cdnjs.cloudflare.com/ajax/libs/bower-angular-translate/2.18.2/angular-translate.min.js \
    653739f29455fa35870da8b0990b176fb23b90fe gui/default/vendor/angular/angular-translate.min.js
get https://cdnjs.cloudflare.com/ajax/libs/bower-angular-translate-loader-static-files/2.18.2/angular-translate-loader-static-files.min.js \
    ac59f8936ec23acd1c70bda29129adc5f9091d77 gui/default/vendor/angular/angular-translate-loader-static-files.min.js

get https://code.jquery.com/jquery-3.5.0.min.js \
    1d6ae46f2ffa213dede37a521b011ec1cd8d1ad3 gui/default/vendor/jquery/jquery.min.js

get https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css \
    b0972fdcce82fd583d4c2ccc3f2e3df7404a19d0 gui/default/vendor/bootstrap/css/bootstrap.min.css
get https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap-theme.min.css \
    beee0e080ea6dcc8c0661b66c1baa08e45f4ecb6 gui/default/vendor/bootstrap/css/bootstrap-theme.min.css
get https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js \
    b16fd8226bd6bfb08e568f1b1d0a21d60247cefb gui/default/vendor/bootstrap/js/bootstrap.min.js

get https://cdnjs.cloudflare.com/ajax/libs/bootstrap-daterangepicker/3.0.5/daterangepicker.min.js \
    cb0e4124f7afc49c49efcf555ae2f88d1203a0f0 gui/default/vendor/daterangepicker/daterangepicker.min.js
get https://cdnjs.cloudflare.com/ajax/libs/bootstrap-daterangepicker/3.0.5/daterangepicker.min.css \
    62a38cf98ecc75c99f9cd0040b46b63a3bd76e4c gui/default/vendor/daterangepicker/daterangepicker.min.css
get https://cdnjs.cloudflare.com/ajax/libs/bootstrap-daterangepicker/3.0.5/moment.min.js \
    69ab16ba8ca68431ab59eff286c7ed1e520bca30 gui/default/vendor/daterangepicker/moment.min.js

get https://cdnjs.cloudflare.com/ajax/libs/jquery.fancytree/2.35.0/jquery.fancytree-all-deps.min.js \
    e4e6acdc86ba20b85692d23416406109afb588dc gui/default/vendor/fancytree/jquery.fancytree-all-deps.min.js
get https://cdnjs.cloudflare.com/ajax/libs/jquery.fancytree/2.35.0/skin-lion/ui.fancytree.min.css \
    9d240c1c35b21afcb2dbc30fa1b9cf5a90206693 gui/default/vendor/fancytree/skin-lion/ui.fancytree.min.css
get https://cdnjs.cloudflare.com/ajax/libs/jquery.fancytree/2.35.0/skin-lion/icons-rtl.gif \
    7c28236756d1400f4a72d46106f18c1dae0a281b gui/default/vendor/fancytree/skin-lion/icons-rtl.gif
get https://cdnjs.cloudflare.com/ajax/libs/jquery.fancytree/2.35.0/skin-lion/icons.gif \
    266100575e2467a3ec6f31b7ec175f11b67a58c3 gui/default/vendor/fancytree/skin-lion/icons.gif
get https://cdnjs.cloudflare.com/ajax/libs/jquery.fancytree/2.35.0/skin-lion/loading.gif \
    dcabdd743fd3e9d7bd5647abeb86e66a3e6f9597 gui/default/vendor/fancytree/skin-lion/loading.gif
get https://cdnjs.cloudflare.com/ajax/libs/jquery.fancytree/2.35.0/skin-lion/vline-rtl.gif \
    d1fd1039ec6c417e59c147bfc1689251371a69fe gui/default/vendor/fancytree/skin-lion/vline-rtl.gif
get https://cdnjs.cloudflare.com/ajax/libs/jquery.fancytree/2.35.0/skin-lion/vline.gif \
    d1fd1039ec6c417e59c147bfc1689251371a69fe gui/default/vendor/fancytree/skin-lion/vline.gif
