pipeline {
  agent {label "predicto-api-"}
    parameters {
        string(name: 'nombre_contenedor', defaultValue: 'predicto-api-container', description: 'nombre del contenedor')
        string(name: 'nombre_imagen', defaultValue: 'predicto-api', description: 'nombre de la imagen')
        string(name: 'tag_imagen', defaultValue: 'latest', description: 'etiqueta de la imagen')
        string(name: 'puerto_imagen', defaultValue: '8081', description: 'puerto de la imagen')
    }
    environment {
        name_final = "${nombre_contenedor}${tag_imagen}${puerto_imagen}"        
    }
    stages {
          stage('stop/rm') {
            when {
                expression { 
                    DOCKER_EXIST = sh(returnStdout: true, script: 'echo "$(docker ps -q --filter name=${name_final})"').trim()
                    return  DOCKER_EXIST != '' 
                }
            }
            steps {
                script{
                    sh ''' 
                        sudo docker stop ${name_final}
                        sudo docker rm -vf ${name_final}
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
                        docker run  -dtp ${puerto_imagen}:${puerto_imagen} --name ${name_final} ${nombre_imagen}:${tag_imagen}
                    '''
                    }
                }                                  
            }
        }
    }