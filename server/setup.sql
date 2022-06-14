create database vault;
create role vault_user with login;
alter user vault_user password 'insecure_password';
GRANT ALL PRIVILEGES ON DATABASE vault TO vault_user;