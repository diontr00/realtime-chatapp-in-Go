if File.exist?('./cert/localhost.cert') && File.exist?('./cert/localhost.key')
  puts 'Certificates already exist.Continute'
else
  # Check if mkcert command is available
  if system('mkcert -h > /dev/null 2>&1')
    puts 'Certificates do not exist, but mkcert is available. Generating certificates...'
    system('mkcert -cert-file ./cert/localhost.cert -key-file ./cert/localhost.key localhost')
    puts 'Certificates generated successfully.'
  else
    puts 'Certificates do not exist, and mkcert is not available. Please install mkcert.'
  end
end
