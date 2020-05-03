(function () {
    'use strict';

    angular
        .module('app')
        .controller('HomeController', ['SourceService', 'FlashService', '$uibModal',
            function (SourceService, FlashService, $uibModal) {
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
                vm.sourceAdd = function() {
                    var modalInstance = $uibModal.open({
                        animation: true,
                        ariaLabelledBy: 'modal-title',
                        ariaDescribedBy: 'modal-body',
                        templateUrl: 'nyansync/modal/modal.add_source.html',
                        controller: 'AddSourceController',
                        controllerAs: 'vm',
                        size: 'lg',
                    });

                    modalInstance.result.then(function () {
                        vm.sources = SourceService.query({cache: false});
                    });
                }
            }
        ]);

})();
