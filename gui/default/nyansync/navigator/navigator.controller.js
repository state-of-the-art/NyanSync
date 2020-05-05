(function () {
  'use strict';

  angular
    .module('app')
    .controller('NavigatorController', ['$scope',
      function ($scope) {
        var vm = this;
        $scope.vm = vm;

        (function initController() {
        })();

        vm.navigator_path = '/asd/qweq/asdqwdq/qweqweqet/qeqweqwe/qwe/asdqwd/qwrqwr/qwe';
        vm.navigator_items = [{
          Name: 'test_file.jpg',
          Preview: '/assets/img/favicon-default.generated.png',
        }];
      }
    ]);

})();
