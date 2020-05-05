(function () {
  'use strict';

  angular
    .module('app')
    .controller('HomeController', [
      function () {
        var vm = this;

        (function initController() {
        })();
      }
    ]).directive('nyansyncsources', function(){
      return {
        restrict: 'A',
        replace: true,
        templateUrl:'nyansync/sources/sources.view.html',
      }
    }).directive('nyansyncnavigator', function(){
      return {
        restrict: 'A',
        replace: true,
        templateUrl:'nyansync/navigator/navigator.view.html',
      }
    });

})();
