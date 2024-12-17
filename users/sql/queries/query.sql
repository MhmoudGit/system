-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (
  username, email, password, first_name, last_name, phone_number, is_active, role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET username = $2,
    email = $3,
    password = $4,
    first_name = $5,
    last_name = $6,
    phone_number = $7,
    is_active = $8,
    role = $9,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ActivateUser :exec
UPDATE users
SET is_active = TRUE,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = FALSE,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;


------------------------ROLES--------------------------------------

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles
WHERE deleted_at IS NULL
ORDER BY role_name;

-- name: CreateRole :one
INSERT INTO roles (
  role_name, permissions
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateRole :exec
UPDATE roles
SET role_name = $2,
    permissions = $3,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteRole :exec
UPDATE roles
SET deleted_at = NOW()
WHERE id = $1;

-- name: HardDeleteRole :exec
DELETE FROM roles
WHERE id = $1;

-- name: AddPermissionToRole :exec
UPDATE roles
SET permissions = array_append(permissions, $2),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RemovePermissionFromRole :exec
UPDATE roles
SET permissions = array_remove(permissions, $2),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

------------------------Permissions------------------------

-- name: GetPermission :one
SELECT * FROM permissions
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListPermissions :many
SELECT * FROM permissions
WHERE deleted_at IS NULL
ORDER BY name;

-- name: CreatePermission :one
INSERT INTO permissions (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdatePermission :exec
UPDATE permissions
SET name = $2,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeletePermission :exec
UPDATE permissions
SET deleted_at = NOW()
WHERE id = $1;

-- name: HardDeletePermission :exec
DELETE FROM permissions
WHERE id = $1;