pipeline {
    agent any
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Test') {
            steps {
                script {
                    docker.image('python:3.9').inside {
                        sh '''
                            sudo pip install pytest
                            sudo pip install -r requirements.txt
                            sudo python -m pytest tests/
                        '''
                    }
                }
            }
        }
    }
}
