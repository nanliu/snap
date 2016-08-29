require "pathname"
require "serverspec"
require "docker"
require "pry"

describe "Dockerfile" do

  before(:all) do
    influxdb_path = File.join Pathname.new(__FILE__).parent.parent, "influxdb/0.12/"
    image = Docker::Image.build_from_dir influxdb_path
    set :backend, :docker
    set :docker_image, image.id
  end

  describe port(8083) do
    it {should be_listening }
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
