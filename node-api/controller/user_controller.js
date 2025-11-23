const UserUseCase = require('../usecase/user_usecase');

// Inicializa uma única instância do caso de uso para ser reutilizada por todos os handlers
const userUseCase = new UserUseCase();

/**
 * Manipula requisições GET /users. Retorna um array JSON com todos os usuários.
 * Em caso de erro, retorna HTTP 500.
 */
exports.getUsers = async (req, res) => {
  try {
    const users = await userUseCase.getUsers();
    res.status(200).json(users);
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: 'erro interno do servidor' });
  }
};

/**
 * Manipula requisições GET /users/:id. Retorna um objeto JSON com o usuário
 * correspondente ao ID fornecido. Se não encontrado, responde com HTTP 404.
 */
exports.getUserById = async (req, res) => {
  const id = parseInt(req.params.id, 10);
  if (isNaN(id)) {
    res.status(400).json({ error: 'id deve ser um inteiro' });
    return;
  }
  try {
    const user = await userUseCase.getUserById(id);
    if (!user) {
      res.status(404).json({ error: 'usuario nao encontrado' });
      return;
    }
    res.status(200).json(user);
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: 'erro interno do servidor' });
  }
};

/**
 * Manipula requisições POST /users. Cria um novo usuário. Em caso de sucesso,
 * retorna HTTP 201 com o usuário criado (sem a senha). Erros de validação e
 * unicidade retornam HTTP 400; outros erros retornam HTTP 500.
 */
exports.createUser = async (req, res) => {
  const user = req.body;
  try {
    const newUser = await userUseCase.createUser(user);
    res.status(201).json(newUser);
  } catch (err) {
    // Erros de validação ou unicidade são retornados como bad request
    if (
      err.message === 'nome é obrigatório' ||
      err.message === 'email é obrigatório' ||
      err.message === 'username é obrigatório' ||
      err.message === 'senha é obrigatória' ||
      err.message === 'email já está em uso' ||
      err.message === 'username já está em uso'
    ) {
      res.status(400).json({ error: err.message });
    } else {
      console.error(err);
      res.status(500).json({ error: 'erro interno do servidor' });
    }
  }
};

/**
 * Manipula requisições PUT /users/:id. Atualiza o usuário especificado. Se o usuário
 * não existir, responde com HTTP 404. Em caso de erro de validação retorna
 * HTTP 400. Caso contrário, retorna o usuário atualizado.
 */
exports.updateUser = async (req, res) => {
  const id = parseInt(req.params.id, 10);
  if (isNaN(id)) {
    res.status(400).json({ error: 'id deve ser um inteiro' });
    return;
  }
  const user = req.body;
  try {
    const updatedUser = await userUseCase.updateUser(id, user);
    res.status(200).json(updatedUser);
  } catch (err) {
    if (err.message === 'usuario nao encontrado') {
      res.status(404).json({ error: 'usuario nao encontrado' });
    } else {
      console.error(err);
      res.status(500).json({ error: 'erro interno do servidor' });
    }
  }
};

/**
 * Manipula requisições DELETE /users/:id. Deleta o usuário especificado. Se não
 * encontrado, responde com HTTP 404. Em caso de sucesso, retorna uma mensagem
 * de confirmação.
 */
exports.deleteUser = async (req, res) => {
  const id = parseInt(req.params.id, 10);
  if (isNaN(id)) {
    res.status(400).json({ error: 'id deve ser um inteiro' });
    return;
  }
  try {
    await userUseCase.deleteUser(id);
    res.status(200).json({ message: 'sucesso ao deletar o usuario' });
  } catch (err) {
    if (err.message === 'usuario nao encontrado') {
      res.status(404).json({ error: 'usuario nao encontrado' });
    } else {
      console.error(err);
      res.status(500).json({ error: 'erro interno do servidor' });
    }
  }
};