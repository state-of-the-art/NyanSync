(function() {
  'use strict';

  angular
    .module('app')
    .controller('LoginController', ['$location', '$localStorage', 'AuthService',
      function( $location, $localStorage, AuthService ) {
        var vm = this;

        (function initController() {
        })();

        vm.login = function login() {
          vm.data_loading = true;
          AuthService.Login(vm.username, vm.password,
            function() {
              vm.data_loading = false;
              $location.path('/');
            }, function() {
              vm.data_loading = false;
            }
          );
        };
      }
    ]);

})();
