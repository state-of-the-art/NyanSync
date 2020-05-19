(function() {
  'use strict';

  angular
    .module('app')
    .controller('SourcesController', ['$scope', 'SourceService', '$uibModal',
      function( $scope, SourceService, $uibModal ) {
        var vm = this;
        $scope.vm = vm;

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
            templateUrl: 'nyanshare/modal/modal.source.html',
            controller: 'SourceController',
            controllerAs: 'vm',
            size: 'lg',
            resolve: {
              source: null,
            },
          }).result.then(function() {
            // Update the whole list of sources from API
            vm.sources = SourceService.query({cache: false});
          });
        };
        vm.sourceEdit = function(source) {
          $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.source.html',
            controller: 'SourceController',
            controllerAs: 'vm',
            size: 'lg',
            resolve: {
              source: function(){ return new SourceService(source); },
            },
          }).result.then(function() {
            // TODO: update cache when URI was changed
            // Update the whole list of sources from API
            vm.sources = SourceService.query({cache: false});
          });
        };
      }
    ]);

})();
