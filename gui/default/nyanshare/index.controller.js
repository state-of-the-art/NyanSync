(function() {
  'use strict';

  angular
    .module('app')
    .controller('IndexController', [ '$rootScope', '$uibModal', 'UserService',
      function( $rootScope, $uibModal, UserService ) {
        var vm = this;

        (function initController() {
        })();

        vm.accountEdit = function() {
          $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.user.html',
            controller: 'UserController',
            controllerAs: 'vm',
            size: 'md',
            resolve: {
              user: function(){ return new UserService($rootScope.account); },
            },
          });
        };
      }
    ]);

})();
