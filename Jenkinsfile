pipeline {
  agent any
    parameters {
        string(name: 'nombre_contenedor', defaultValue: 'predicto_api_container', description: 'nombre del contenedor')
        string(name: 'nombre_imagen', defaultValue: 'predicto_api', description: 'nombre de la imagen')
        string(name: 'tag_imagen', defaultValue: 'latest', description: 'etiqueta de la imagen')
        string(name: 'puerto_imagen', defaultValue: '8081', description: 'puerto de la imagen')
    }
    environment {
        nombre_final = "${nombre_contenedor}${tag_imagen}${puerto_imagen}"        
    }
    stages {
          stage('stop/rm') {
            when {
                expression { 
                    DOCKER_EXIST = sh(returnStdout: true, script: 'echo "$(docker ps -q --filter name=${nombre_final})"').trim()
                    return  DOCKER_EXIST != '' 
                }
            }
            steps {
                script{
                    sh ''' 
                        sudo docker stop ${env.nombre_final}
                        sudo docker rm -vf ${env.nombre_final}
                    '''
                    }
                }                                   
            }
        stage('build') {
            steps {
                script{
                    sh ''' 
                    docker build . -t ${nombre_imagen}:${tag_imagen}
                    '''
                    }
                }                                       
            }
            stage('run') {
            steps {
                script{
                    sh ''' 
                        docker run  -dtp ${puerto_imagen}:${puerto_imagen} --name ${nombre_final} ${nombre_imagen}:${tag_imagen}
                    '''
                    }
                }                                  
            }
        }
    }