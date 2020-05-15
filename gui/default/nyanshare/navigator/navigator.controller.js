(function() {
  'use strict';

  angular
    .module('app')
    .controller('NavigatorController', ['$scope', '$uibModal', 'NavigatorService', 'DropdownService',
      function( $scope, $uibModal, NavigatorService, DropdownService ) {
        var vm = this;
        $scope.vm = vm;

        (function initController() {
        })();

        vm.navigator_path = [];
        vm.navigator_active = null;
        vm.navigator_active_el = null;

        vm.navigatePath = function(path) {
          vm.navigator_path = path;
          vm.navigator_items = NavigatorService.query(vm.navigator_path);
        };

        vm.navigateChildren = function(name) {
          vm.navigatePath(vm.navigator_path.concat([name]));
        };

        vm.itemMenu = function(item, e) {
          if( vm.navigator_active && vm.navigator_active !== item ) {
            DropdownService.Remove($(vm.navigator_active_el).parent());
          }
          vm.navigator_active = item;
          vm.navigator_active_el = e.currentTarget;
          var p = $(e.currentTarget).parent();
          if( DropdownService.Exist(p) )
            return;

          // Create dropdown menu
          DropdownService.Create({
            'Share': function() {
              vm.shareItem(item);
            },
          }, p);
        };

        vm.itemClick = function(item, e) {
          if( item === vm.navigator_active ) {
            if( item.Type == 'folder' )
              vm.navigateChildren(item.Name);
          } else {
            // For touch devices to use double click system
            if( vm.navigator_active_el ) {
              $(vm.navigator_active_el)
                .trigger('mouseout')
                .removeClass('active');
            }
            $(e.currentTarget)
              .trigger('mouseover')
              .addClass('active');

            vm.itemMenu(item, e);
          }
        };

        vm.shareItem = function(item) {
          $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.access.html',
            controller: 'AccessController',
            controllerAs: 'vm',
            size: 'lg',
            resolve: {
              title: function(){ return 'Create new share'; },
              source_id: function(){ return vm.navigator_path[0] || item.Name; },
              path: function(){
                if( vm.navigator_path.length == 0 )
                  return '';
                else
                  return vm.navigator_path.slice(1).concat(item.Name).join('/');
              },
              item: function(){ return item; },
            },
          }).result.then(function() {
            // Update the whole list of sources from API
            vm.sources = SourceService.query({cache: false});
          });
        };

        vm.navigatePath([]);
      }
    ]);

})();
