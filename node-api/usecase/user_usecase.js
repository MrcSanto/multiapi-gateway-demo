const bcrypt = require('bcrypt');
const UserRepository = require('../repository/user_repository');

/**
 * UserUseCase encapsulates business rules for managing users. It relies on
 * UserRepository to interact with the database and handles validation and
 * password hashing logic before delegating to the repository layer.
 */
class UserUseCase {
  constructor() {
    this.repository = new UserRepository();
  }

  /**
   * Fetches all users.
   *
   * @returns {Promise<Array>}
   */
  async getUsers() {
    return this.repository.getUsers();
  }

  /**
   * Retrieves a single user by ID.
   *
   * @param {number} id - User identifier.
   * @returns {Promise<Object|null>}
   */
  async getUserById(id) {
    return this.repository.getUserById(id);
  }

  /**
   * Deletes a user by ID. Throws an error if no user was removed.
   *
   * @param {number} id - User identifier to remove.
   * @returns {Promise<void>}
   */
  async deleteUser(id) {
    const rowCount = await this.repository.deleteUser(id);
    if (rowCount === 0) {
      throw new Error('usuario nao encontrado');
    }
  }

  /**
   * Updates an existing user. If the user does not exist, an error is thrown.
   *
   * @param {number} id - User identifier.
   * @param {Object} user - Data to update.
   * @returns {Promise<Object>} The updated user record.
   */
  async updateUser(id, user) {
    const rowCount = await this.repository.updateUser(id, user);
    if (rowCount === 0) {
      throw new Error('usuario nao encontrado');
    }
    // Return the updated user
    return this.repository.getUserById(id);
  }

  /**
   * Creates a new user with validation and password hashing. Throws
   * descriptive errors if validation fails or uniqueness constraints are
   * violated.
   *
   * @param {Object} user - Data for the new user.
   * @returns {Promise<Object>} The newly created user with the ID field set.
   */
  async createUser(user) {
    // Basic validations
    if (!user.name) {
      throw new Error('nome é obrigatório');
    }
    if (!user.email) {
      throw new Error('email é obrigatório');
    }
    if (!user.username) {
      throw new Error('username é obrigatório');
    }
    if (!user.password) {
      throw new Error('senha é obrigatória');
    }

    // Check if email or username already exists
    const existingByEmail = await this.repository.getUserByEmail(user.email);
    if (existingByEmail) {
      throw new Error('email já está em uso');
    }
    const existingByUsername = await this.repository.getUserByUserName(user.username);
    if (existingByUsername) {
      throw new Error('username já está em uso');
    }

    // Hash the password using bcrypt
    const hashed = await bcrypt.hash(user.password, 10);
    const newUser = {
      name: user.name,
      email: user.email,
      username: user.username,
      password: hashed,
    };

    // Insert into the database
    const id = await this.repository.createUser(newUser);
    return {
      id,
      name: newUser.name,
      email: newUser.email,
      username: newUser.username,
    };
  }
}

module.exports = UserUseCase;