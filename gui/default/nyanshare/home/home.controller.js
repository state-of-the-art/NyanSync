(function() {
  'use strict';

  angular
    .module('app')
    .controller('HomeController', [
      function() {
        var vm = this;

        (function initController() {
        })();
      }
    ]).directive('nyanshareSources', function(){
      return {
        restrict: 'A',
        replace: true,
        templateUrl:'nyanshare/sources/sources.view.html',
      }
    }).directive('nyanshareNavigator', function(){
      return {
        restrict: 'A',
        replace: true,
        templateUrl:'nyanshare/navigator/navigator.view.html',
      }
    });

})();
