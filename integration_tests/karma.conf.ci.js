const getConf = require('./karma.conf.common')
module.exports = function(config) {
  config.set(getConf(true, config))
}
