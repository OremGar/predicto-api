pipeline {
  agent any
    parameters {
        string(name: 'nombre_contenedor', defaultValue: 'predicto_api_container', description: 'nombre del contenedor')
        string(name: 'nombre_imagen', defaultValue: 'predicto_api', description: 'nombre de la imagen')
        string(name: 'tag_imagen', defaultValue: 'latest', description: 'etiqueta de la imagen')
        string(name: 'puerto_imagen', defaultValue: '8081', description: 'puerto de la imagen')
        string(name: 'puerto_externo', defaultValue: '8083', description: 'puerto externo')
    }
    environment {
        nombre_final = "${nombre_contenedor}"        
    }
    stages {
          stage('stop/rm') {
            when {
                expression { 
                    DOCKER_EXIST = sh(returnStdout: true, script: 'echo "$(docker ps -a -q --filter name=${nombre_final})"').trim()
                    return  DOCKER_EXIST != '' 
                }
            }
            steps {
                script{
                    sh """
                        docker stop ${nombre_final}
                        docker rm -vf ${nombre_final}
                    """
                    }
                }                                   
            }
        stage('build') {
            steps {
                script{
                    sh """
                    docker build . -t ${nombre_imagen}:${tag_imagen}
                    """
                    }
                }                                       
            }
            stage('run') {
            steps {
                script{
                    sh """ 
                        docker run  -dtp ${puerto_externo}:${puerto_imagen} --name ${nombre_final} ${nombre_imagen}:${tag_imagen}
                    """
                    }
                }                                  
            }
        }
    }