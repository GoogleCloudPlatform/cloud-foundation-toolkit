/**
 * Responds to any HTTP request with 'Hello World!'.
 *
 * @param {Object} req Cloud Function request context.
 * @param {Object} res Cloud Function response context.
 */
exports.helloGET = function (req, res) {
  res.status(200).send('Hello world!');
};

