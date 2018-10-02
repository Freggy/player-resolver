Vagrant.configure('2') do |config|
	dataDir = File.join(File.dirname(__FILE__), '.vagrant/cassandra_data')
	config.vm.provider "docker" do |d|
		d.image   = "cassandra:latest"
		d.volumes = ["#{dataDir}:/var/lib/cassandra"]
		d.ports   = ["7000:7000", "7001:7001", "7199:7199", "9042:9042", "9160:9160", "9142:9142"]
		d.name    = "pr-local"
	end
end
