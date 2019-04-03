#!/usr/bin/env groovy


// requires "Pipeline Utility Steps" plugin
// requires "pipeline-vars" file to be setup inside jenkins user home dir

def config_dir = "pipeline/network pipeline/app" // relative to ${cft_dir}

def env = "~/pipeline-vars"
def repo = "https://github.com/sourced/deploymentmanager-samples"
def branch = "pipeline"
def git_dir = "cft"
def cft_dir = "community/cloud-foundation"

pipeline {
    agent any
    stages {
        stage('Checkout Repos') {
            steps {
                sh "rm -rf ${git_dir}"
                sh "git clone ${repo} ${git_dir}"
                sh "cd cft && git checkout ${branch}"
            }
        }
        stage('Initialize Deploy Stages') {
            steps {
                sh ". ${env} && cd ${git_dir}/${cft_dir} && cft delete ${config_dir} --show-stages --format yaml> .stages.yaml"
                script {
                    def graph = readYaml file: "${git_dir}/${cft_dir}/.stages.yaml"
                    def i = 1
                    graph.each { stg ->
                        stage("stage-${i}") {
                            def config_list = []
                            stg.each { conf ->
                                config_list.add(conf.source)
                            }
                            def configs = config_list.join(" ")
                            echo "Executing configs: ${configs}"
                            sh ". ${env} && cd ${git_dir}/${cft_dir} && cft delete ${configs}"
                        }
                        i++
                    }
                }
            }
        }
    }
}
