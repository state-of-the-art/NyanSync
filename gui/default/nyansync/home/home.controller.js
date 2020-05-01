(function () {
    'use strict';

    angular
        .module('app')
        .controller('HomeController', ['SourceService', 'FlashService',
            function (SourceService, FlashService) {
                var vm = this;

                (function initController() {
                })();

                vm.sources = SourceService.query();
                vm.sourcesShares = function(source) {
                    return ""
                }
                vm.sourceSetPause = function(source_id, val) {}
                vm.allSourcesSetPause = function(val) {}
                vm.sourceIsAtleastOnePausedStateSetTo = function(val) { return true }
                vm.sourceAdd = function() {}
            }
        ]);

})();
