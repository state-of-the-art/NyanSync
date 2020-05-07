(function () {
  'use strict';

  angular
    .module('app')
    .controller('NavigatorController', ['$scope', 'NavigatorService',
      function ($scope, NavigatorService) {
        var vm = this;
        $scope.vm = vm;

        (function initController() {
        })();

        vm.navigator_path = [];

        vm.navigatePath = function(path) {
          vm.navigator_path = path;
          vm.navigator_items = NavigatorService.query(vm.navigator_path);
        }
        vm.navigateChildren = function(name) {
          vm.navigatePath(vm.navigator_path.concat([name]));
        }
        vm.itemClick = function(item) {
          vm.navigateChildren(item.Name);
        };

        vm.navigatePath([]);
      }
    ]);

})();
