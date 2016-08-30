require "pathname"
require "serverspec"
require "docker"
require "pry"

RSpec.configure do |config|
  config.docker_container_create_options = {
    'HostConfig': {
      "PortBindings": {
        "8083/tcp": [{ 'HostPort' => '8083' }],
        "8086/tcp": [{ 'HostPort' => '8086' }],
      },
    }
  }
end
describe "Dockerfile" do

  before(:all) do
    influxdb_path = File.join Pathname.new(__FILE__).parent.parent, "influxdb/0.12/"

    image = Docker::Image.get "7bb212dd9198"
    set :backend, :docker
    set :docker_image, image.id
  end

  describe port(8083) do
    it {
      should be_listening
    }
  end

  # TODO: These wont' work
  #describe command ('SHOW DATABASES') do
  #  its(:stdout){should contain}
  #  ('playground')
  #end

  #describe command ('SHOW DIAGNOSTICS') do
  #  its(:stdout){should contain}
  #  ('0.13.0')
  #end
end
