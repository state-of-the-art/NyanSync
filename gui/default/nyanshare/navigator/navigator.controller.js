(function() {
  'use strict';

  angular
    .module('app')
    .controller('NavigatorController', ['$scope', '$uibModal', '$localStorage', 'urls', 'NavigatorService', 'SourceService', 'DropdownService',
      function( $scope, $uibModal, $localStorage, urls, NavigatorService, SourceService, DropdownService ) {
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
          var menu = {
            'Share': function() {
              vm.shareItem(item);
              return false;
            },
          };
          if( item.Type == 'source') {
            menu['Edit'] = function() {
              SourceService.get({'Id': item.Name}).$promise.then(function(source){
                vm.sourceEdit(source);
              });
              return false;
            };
          }
          if( item.Type == 'binary') {
            menu['Download'] = function() {
              /* This way downloads file without asking user to save
              $http.get(urls.BASE_API + 'download/' + vm.navigator_path.concat(item.Name).join('/'), {
                responseType: 'blob',
              });*/
              // TODO: replace with a good download way using one-time token
              window.location = urls.BASE_API + 'download/' + vm.navigator_path.concat(item.Name).join('/') +
                '?token=' + $localStorage.account.token;
              return false;
            };
          }
          DropdownService.Create(menu, p);
        };

        vm.itemClick = function(item, e) {
          if( item === vm.navigator_active ) {
            if( ['folder', 'source'].includes(item.Type) )
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
              access: function(){
                return {
                  SourceId: vm.navigator_path[0] || item.Name,
                  Path: vm.navigator_path.length == 0 ? '' : vm.navigator_path.slice(1).concat(item.Name).join('/'),
                }
              },
            },
          }).result.then(function() {
            // Update navigator view
            vm.navigatePath(vm.navigator_path);
          });
        };

        vm.navigatePath([]);

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
            // Update navigator view
            vm.navigatePath(vm.navigator_path);
          });
        };
        vm.sourceEdit = function(source) {
          console.log("Run source edit with ", source);
          $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.source.html',
            controller: 'SourceController',
            controllerAs: 'vm',
            size: 'lg',
            resolve: {
              source: function(){ return source; },
            },
          }).result.then(function() {
            // Update navigator view
            vm.navigatePath(vm.navigator_path);
          });
        };
      }
    ]);

})();
