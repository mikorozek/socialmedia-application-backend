pipeline {
    agent any

    environment {
        NEXUS_URL = 'http://52.233.173.205:8081'
        NEXUS_REPO = 'maven-repository'
    }

    stages {
        stage('Upload to Nexus') {
            agent {
                docker {
                    image 'maven:3.8.4-openjdk-11'
                    args '-v $HOME/.m2:/root/.m2'
                }
            }
            steps {
                withCredentials([usernamePassword(credentialsId: 'nexus-admin', usernameVariable: 'NEXUS_USERNAME', passwordVariable: 'NEXUS_PASSWORD')]) {
                    script {
                        // Tworzymy dynamicznie plik settings.xml
                        sh """
                        mkdir -p ~/.m2
                        cat > ~/.m2/settings.xml <<EOF
                        <settings>
                          <servers>
                            <server>
                              <id>nexus</id>
                              <username>${NEXUS_USERNAME}</username>
                              <password>${NEXUS_PASSWORD}</password>
                            </server>
                          </servers>
                        </settings>
                        EOF
                        """
                        // Uruchamiamy Maven deploy
                        sh """
                        mvn deploy -DskipTests \
                            -DaltDeploymentRepository=nexus::default::${NEXUS_URL}/repository/${NEXUS_REPO}
                        """
                    }
                }
            }
        }
    }
}
