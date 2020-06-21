(function() {
  'use strict';

  angular
    .module('app')
    .controller('UserController', ['user', '$scope', '$uibModalInstance', '$uibModal', 'UserService', 'AuthService', 'RoleService',
      function( user, $scope, $uibModalInstance, $uibModal, UserService, AuthService, RoleService ) {
        var vm = this;

        if( user instanceof UserService )
          vm.title = 'Edit user "' + user.Name + '"';
        else
          vm.title = 'Create new user';

        (function initController() {
          // User to create or edit
          $scope.user = user instanceof UserService ? user : new UserService(user);

          // Set original id for renaming
          $scope.user._orig_login = $scope.user.Login;
          if( !$scope.user._orig_login )
            $scope.user.Login = $scope.user.Name;

          // Check manager is set - and if not - set the current user
          if( $scope.user.Manager === undefined ) {
            $scope.user.Manager = AuthService.GetTokenClaims().id;
          }

          $scope.user_roles = [];
          if( $scope.user.Roles !== null ) {
            for( var r of $scope.user.Roles )
              $scope.user_roles.push({Id: r});
          }

          $scope.used_logins = [];
          // To validate the used user logins
          UserService.query().$promise.then(function(users){
            // Make sure while edit we can reuse the same login of user
            $scope.used_logins = users.map(u => u.Login != $scope.user.Login ? u.Login : null);
          });
        })();

        vm.submit = function() {
          if( !$scope.user._orig_login )
            $scope.user._orig_login = $scope.user.Login

          // Add list of roles to the user object
          $scope.user.Roles = [];
          for( const r of $scope.user_roles )
            $scope.user.Roles.push(r.Id);

          $scope.user.$save().then(function() {
            $uibModalInstance.close($scope.user);
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
                return 'Are you sure you want to remove user "' + $scope.user.Name +
                  '" with login "' + $scope.user._orig_login + '"?';
              },
            },
          }).result.then(function( result ) {
            if( result === true ) {
              $scope.user.$remove().then(function(){
                $uibModalInstance.close();
              });
            }
          });
        };

        vm.loadRoles = function(query) {
          if( query )
            return RoleService.query({q: query, cache: false}).$promise;
          return RoleService.query().$promise;
        };
        vm.checkRole = function(data) {
          return RoleService.get({Id: data.Id}).$promise.then(function(role) {
            // Check the role is present only once
            for( const r of $scope.user_roles ) {
              if( r.Id === role.Id )
                return false;
            }
            Object.assign(data, role);
          }, function(error) {
            return vm.createRole(data);
          });
        };
        vm.createRole = function(data) {
          var promise = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: 'nyanshare/modal/modal.role.html',
            controller: 'RoleController',
            controllerAs: 'vm',
            size: 'md',
            resolve: {
              role: function(){ return data; },
            },
          }).result.then(function( result ) {
            RoleService.query({cache: false});
            if( result instanceof RoleService ) {
              data.Id = result.Id
            }
          });
          return promise;
        };
      }
    ])
    .directive('myUserUnique', function() {
      return {
        require: 'ngModel',
        scope: true,
        link: function(scope, elm, attrs, ctrl) {
          if( !ctrl )
            return;
          ctrl.$validators.userunique = function(modelValue, viewValue) {
            return scope.used_logins.indexOf(modelValue) === -1;
          };
        },
      };
    });

})();
