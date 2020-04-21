angular.module('nyansync.core')
    .filter('metric', function () {
        return function (input) {
            return unitPrefixed(input, false);
        };
    });
