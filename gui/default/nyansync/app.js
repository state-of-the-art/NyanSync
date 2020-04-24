(function () {
    'use strict';

    angular.module('app', [
        'ngStorage',
        'ngRoute',
    ])
		.constant('urls', {
            BASE: '/',
            BASE_API: '/api/v1'
        })
        .config(['$routeProvider', '$httpProvider', function($routeProvider, $httpProvider) {
            $routeProvider
                .when('/', {
                    controller: 'HomeController',
                    templateUrl: 'nyansync/home/home.view.html',
                    controllerAs: 'vm',
                })
                .when('/login', {
                    controller: 'LoginController',
                    templateUrl: 'nyansync/login/login.view.html',
                    controllerAs: 'vm',
                })
                .otherwise({ redirectTo: '/' });

            $httpProvider.interceptors.push(['$q', '$location', '$localStorage', function ($q, $location, $localStorage) {
                return {
                    'request': function (config) {
                        config.headers = config.headers || {};
                        if ($localStorage.token) {
                            config.headers.Authorization = 'Bearer ' + $localStorage.token;
                        }
                        return config;
                    },
                    'responseError': function (response) {
                        if (response.status === 401 || response.status === 403) {
                            delete $localStorage.token;
                            $location.path('/login');
                        }
                        return $q.reject(response);
                    }
                };
            }]);
        }])
        .run(function($rootScope, $location, $localStorage) {
            $rootScope.$on( "$routeChangeStart", function(event, next) {
                if ($localStorage.token == null) {
                    $location.path("/login");
                }
            });
        });

})();
