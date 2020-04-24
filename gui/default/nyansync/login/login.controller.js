(function () {
    'use strict';

    angular
        .module('app')
        .controller('LoginController', ['$location', '$localStorage', 'AuthService', 'FlashService',
            function ($location, $localStorage, AuthService, FlashService) {
                var vm = this;

                vm.login = login;

                (function initController() {
                    // reset login status
                    //AuthService.ClearCredentials();
                })();

                function login() {
                    vm.dataLoading = true;
                    AuthService.Login(vm.username, vm.password,
                        function (res) {
                            $localStorage.token = res.data.token;
                            vm.dataLoading = false;
                            $location.path('/');
                        }, function (res) {
                            FlashService.Error(res.data.message);
                            vm.dataLoading = false;
                        }
                    );
                };
            }
        ]);

})();
