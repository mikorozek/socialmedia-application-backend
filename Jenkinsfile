pipeline {
    agent any

    environment {
        NEXUS_URL = 'http://52.233.173.205:8081'  // Adres Twojego Nexusa
        NEXUS_REPO = 'maven-repository'  // Repozytorium w Nexusie, gdzie będziesz wysyłać artefakty
        NEXUS_CREDENTIALS = 'nexus-admin2'  // ID poświadczeń w Jenkinsie
        MAVEN_HOME = '/usr/share/maven'   // Ścieżka do Mavena w kontenerze Docker (jeśli potrzebne)
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }
            steps {
                script {
                    sh 'mvn clean package'
                }
            }
        }

        stage('Test') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }
            steps {
                script {
                    sh 'mvn test'
                }
            }
        }

        stage('Upload to Nexus') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'nexus-admin', usernameVariable: 'NEXUS_USERNAME', passwordVariable: 'NEXUS_PASSWORD')]) {
                        sh """
                            if [ ! -f ~/.m2/settings.xml ]; then
                                mkdir -p ~/.m2
                                echo '<settings></settings>' > ~/.m2/settings.xml
                            fi
                            sed -i '/<servers>/,/<\\/servers>/d' ~/.m2/settings.xml
                            sed -i '/<\\/settings>/i\\
                            <servers>\\
                                <server>\\
                                    <id>nexus</id>\\
                                    <username>${NEXUS_USERNAME}</username>\\
                                    <password>${NEXUS_PASSWORD}</password>\\
                                </server>\\
                            </servers>' ~/.m2/settings.xml
                        """
                    }

                    sh "mvn deploy -DskipTests"
                }
}

        }
    }
}
