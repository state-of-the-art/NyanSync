(function() {
  'use strict';

  angular
    .module('app')
    .controller('AccessController', ['access', '$scope', '$uibModalInstance', '$uibModal', 'AccessService', 'UserService', 'SourceService', 'AuthService',
      function( access, $scope, $uibModalInstance, $uibModal, AccessService, UserService, SourceService, AuthService ) {
        var vm = this;

        if( access instanceof AccessService )
          vm.title = 'Edit access "' + access.Id + '"';
        else
          vm.title = 'Create new access';

        (function initController() {
          // Access to create or edit
          $scope.access = access instanceof AccessService ? access : new AccessService(access);

          // Check manager is set - and if not - set the current user
          if( $scope.access.Manager === undefined ) {
            $scope.access.Manager = AuthService.GetTokenClaims().id;
          }

          $scope.access_users = [];
          if( $scope.access.Users !== null ) {
            // TODO: fill access_users with actual objects
          }

          // To validate the source ids
          $scope.source_ids = []
          SourceService.query().$promise.then(function(sources){
            $scope.source_ids = sources.map(s => s.Id);
          });
        })();

        vm.submit = function() {
          // Add list of user logins to access object
          $scope.access.Users = [];
          for( const u of $scope.access_users )
            $scope.access.Users.push(u.Login);
          $scope.access.$save().then(function(){
            $uibModalInstance.close();
          });
        };

        vm.cancel = function() {
          $uibModalInstance.dismiss('cancel');
        };

        vm.remove = function() {
          $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.question.html',
            controller: 'QuestionController',
            controllerAs: 'vm',
            size: 'sm',
            resolve: {
              body: function(){
                return 'Are you sure you want to remove access "' + $scope.access.Id + '"?';
              },
            },
          }).result.then(function( result ) {
            if( result === true ) {
              $scope.access.$remove().then(function(){
                $uibModalInstance.close();
              });
            }
          });
        };

        vm.loadUsers = function(query) {
          if( query )
            return UserService.query({q: query, cache: false}).$promise;
          return UserService.query().$promise;
        };
        vm.checkUser = function(data) {
          return UserService.get({Login: data.Login || data.Name}).$promise.then(function(user) {
            // Check the user is present only once
            for( const u of $scope.access_users ) {
              if( u.Login === user.Login )
                return false;
            }
            Object.assign(data, user);
          }, function(error) {
            return vm.createUser(data);
          });
        };
        vm.createUser = function(data) {
          // Set guest role by default
          data.Roles = ['guest'];
          var promise = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.user.html',
            controller: 'UserController',
            controllerAs: 'vm',
            size: 'md',
            resolve: {
              user: function(){ return data; },
            },
          }).result.then(function( result ) {
            UserService.query({cache: false});
            if( result instanceof UserService ) {
              data.Login = result.Login
              data.Name = result.Name
            }
          });
          return promise;
        };
      }
    ])
    .directive('mySourceIdExist', function() {
      return {
        require: 'ngModel',
        scope: true,
        link: function(scope, elm, attrs, ctrl) {
          if( !ctrl )
            return;
          ctrl.$validators.sourceid = function(modelValue, viewValue) {
            return scope.source_ids.indexOf(modelValue) !== -1;
          };
        },
      };
    });

})();
