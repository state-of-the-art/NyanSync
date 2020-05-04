(function () {
    'use strict';

    angular
        .module('app')
        .controller('HomeController', ['SourceService', '$uibModal',
            function (SourceService, $uibModal) {
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
                    $uibModal.open({
                        animation: true,
                        ariaLabelledBy: 'modal-title',
                        ariaDescribedBy: 'modal-body',
                        templateUrl: 'nyansync/modal/modal.source.html',
                        controller: 'SourceController',
                        controllerAs: 'vm',
                        size: 'lg',
                        resolve: {
                            title: function(){ return 'Add new source'; },
                            source: null,
                        },
                    }).result.then(function () {
                        // Update the whole list of sources from API
                        vm.sources = SourceService.query({cache: false});
                    });
                }
                vm.sourceEdit = function(source) {
                    $uibModal.open({
                        animation: true,
                        ariaLabelledBy: 'modal-title',
                        ariaDescribedBy: 'modal-body',
                        templateUrl: 'nyansync/modal/modal.source.html',
                        controller: 'SourceController',
                        controllerAs: 'vm',
                        size: 'lg',
                        resolve: {
                            title: function(){ return 'Edit source'; },
                            source: function(){ return source; },
                        },
                    }).result.then(function () {
                        // Update the whole list of sources from API
                        vm.sources = SourceService.query({cache: false});
                    });
                }
            }
        ]);

})();
