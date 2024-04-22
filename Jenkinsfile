pipeline {
  agent any
    parameters {
        string(name: 'nombre_contenedor', defaultValue: 'predicto_api_container', description: 'nombre del contenedor')
        string(name: 'nombre_imagen', defaultValue: 'predicto_api', description: 'nombre de la imagen')
        string(name: 'tag_imagen', defaultValue: 'latest', description: 'etiqueta de la imagen')
        string(name: 'puerto_imagen', defaultValue: '8081', description: 'puerto de la imagen')
<<<<<<< HEAD
    }
    environment {
        name_final = "${nombre_contenedor}${tag_imagen}${puerto_imagen}"        
=======
        string(name: 'puerto_externo', defaultValue: '8083', description: 'puerto externo')
    }
    environment {
        nombre_final = "${nombre_contenedor}"        
>>>>>>> dev
    }
    stages {
          stage('stop/rm') {
            when {
                expression { 
<<<<<<< HEAD
                    DOCKER_EXIST = sh(returnStdout: true, script: 'echo "$(docker ps -q --filter name=${name_final})"').trim()
=======
                    DOCKER_EXIST = sh(returnStdout: true, script: 'echo "$(docker ps -a -q --filter name=${nombre_final})"').trim()
>>>>>>> dev
                    return  DOCKER_EXIST != '' 
                }
            }
            steps {
                script{
<<<<<<< HEAD
                    sh ''' 
                        sudo docker stop ${name_final}
                        sudo docker rm -vf ${name_final}
                    '''
=======
                    sh """
                        docker stop ${nombre_final}
                        docker rm -vf ${nombre_final}
                    """
>>>>>>> dev
                    }
                }                                   
            }
        stage('build') {
            steps {
                script{
<<<<<<< HEAD
                    sh ''' 
                    docker build . -t ${nombre_imagen}:${tag_imagen}
                    '''
=======
                    sh """
                    docker build . -t ${nombre_imagen}:${tag_imagen}
                    """
>>>>>>> dev
                    }
                }                                       
            }
            stage('run') {
            steps {
                script{
<<<<<<< HEAD
                    sh ''' 
                        docker run  -dtp ${puerto_imagen}:${puerto_imagen} --name ${name_final} ${nombre_imagen}:${tag_imagen}
                    '''
=======
                    sh """ 
                        docker run  -dtp ${puerto_externo}:${puerto_imagen} --name ${nombre_final} ${nombre_imagen}:${tag_imagen}
                    """
>>>>>>> dev
                    }
                }                                  
            }
        }
    }