(function () {
    'use strict';

    angular
        .module('app')
        .controller('LoginController', ['$location', '$localStorage', 'AuthService', 'FlashService',
            function ($location, $localStorage, AuthService, FlashService) {
                var vm = this;

                (function initController() {
                })();

                vm.login = function login() {
                    vm.dataLoading = true;
                    AuthService.Login(vm.username, vm.password,
                        function (res) {
                            $localStorage.token = res.data.token;
                            vm.dataLoading = false;
                            $location.path('/');
                        }, function (res) {
                            vm.dataLoading = false;
                        }
                    );
                };
            }
        ]);

})();
