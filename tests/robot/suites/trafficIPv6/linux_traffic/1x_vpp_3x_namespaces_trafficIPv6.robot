*** Settings ***
Library      OperatingSystem
#Library      RequestsLibrary
#Library      SSHLibrary      timeout=60s
#Library      String

Resource     ../../../variables/${VARIABLES}_variables.robot

Resource     ../../../libraries/all_libs.robot

Force Tags        traffic     IPv6    ExpectedFailure
Suite Setup       Testsuite Setup
Suite Teardown    Testsuite Teardown
Test Setup        TestSetup
Test Teardown     TestTeardown
*** Variables ***
${VARIABLES}=          common
${ENV}=                common
${WAIT_TIMEOUT}=     20s
${SYNC_SLEEP}=       3s
${RESYNC_SLEEP}=       1s
# wait for resync vpps after restart
${RESYNC_WAIT}=        30s
${IP_1}=         fd31::1:a:0:0:1
${IP_2}=         fd31::1:a:0:0:2
${IP_3}=         fd31::1:a:0:0:5
${IP_4}=         fd31::1:a:0:0:6

*** Test Cases ***
Configure Environment
    [Tags]    setup
    Add Agent VPP Node    agent_vpp_1

Show Interfaces Before Setup
    vpp_term: Show Interfaces    agent_vpp_1

Setup Interfaces
    Put Veth Interface Via Linux Plugin    node=agent_vpp_1    namespace=ns1    name=ns1_veth1    host_if_name=ns1_veth1_linux    mac=d2:74:8c:12:67:d2    peer=ns2_veth2    ip=192.168.22.1    prefix=30
    Put Veth Interface Via Linux Plugin    node=agent_vpp_1    namespace=ns2    name=ns2_veth2    host_if_name=ns2_veth2_linux    mac=92:c7:42:67:ab:cd    peer=ns1_veth1    ip=${IP_2}    prefix=30

    Put Veth Interface Via Linux Plugin    node=agent_vpp_1    namespace=ns2    name=ns2_veth3    host_if_name=ns2_veth3_linux    mac=92:c7:42:67:ab:cf    peer=ns3_veth3    ip=${IP_3}    prefix=30
    Put Veth Interface Via Linux Plugin    node=agent_vpp_1    namespace=ns3    name=ns3_veth3    host_if_name=ns3_veth3_linux    mac=92:c7:42:67:ab:ce    peer=ns2_veth3    ip=${IP_4}    prefix=30



Chcek Linux Interfaces
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Interfaces    node=agent_vpp_1    namespace=ns1    interface=ns1_veth1
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Interfaces    node=agent_vpp_1    namespace=ns2    interface=ns2_veth2
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Interfaces    node=agent_vpp_1    namespace=ns2    interface=ns2_veth3
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Interfaces    node=agent_vpp_1    namespace=ns3    interface=ns3_veth3

Ping6 In Namespaces
    # This should work by default after veth interface setup
    Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_2}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns2    ip=${IP_1}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns2    ip=${IP_4}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_3}

# This will fail now
#     Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_3}
#     Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_4}
#
#     Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_1}
#     Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_2}

Create Linux Defalut Routes
    # this did not work
    #Put Linux Route    node=agent_vpp_1    namespace=ns1    interface=ns1_veth1    routename=innercross1    ip=${IP_3}    prefix=32    next_hop=${IP_1}
    #Put Linux Route    node=agent_vpp_1    namespace=ns3    interface=ns3_veth3    routename=innercross2    ip=${IP_2}    prefix=32    next_hop=${IP_4}
    Put Default Linux Route    node=agent_vpp_1    namespace=ns1    interface=ns1_veth1    routename=innercross1    next_hop=${IP_1}
    Put Default Linux Route    node=agent_vpp_1    namespace=ns3    interface=ns3_veth3    routename=innercross2    next_hop=${IP_4}
    Sleep    ${CONFIG_SLEEP}

Check Linux Default Routes
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Default Routes    node=agent_vpp_1    namespace=ns1    next_hop=${IP_1}
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Default Routes    node=agent_vpp_1    namespace=ns3    next_hop=${IP_4}

Ping6 In Namespaces Again
    Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_3}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_2}

# This will fail now
#     Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_4}
#     Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_1}

Create Linux Routes2
    #This needs to be fixed  - https://jira.pantheon.sk/browse/ODPM-743
    Put Linux Route Without Interface    node=agent_vpp_1    namespace=ns1    routename=outercross1    ip=${IP_4}    prefix=32    next_hop=${IP_2}
    Put Linux Route Without Interface    node=agent_vpp_1    namespace=ns3    routename=outercross2    ip=${IP_1}    prefix=32    next_hop=${IP_3}


    #temporarily - because previous commands does not work
    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns1 ip route add ${IP_4}/32 via ${IP_2}
    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns3 ip route add ${IP_1}/32 via ${IP_3}

    # IP FORWARDING
    # Advice from Andrej Marcinek
    # 9, enable IP forwarding
    # sudo sysctl -w net.ipv4.ip_forward=1
    # 10, do 'sudo nano /etc/sysctl.conf' and uncomment line:
    # net.ipv4.ip_forward=1
    #https://unix.stackexchange.com/questions/292801/routing-between-linux-namespaces
    #ip netns exec $NS_MID sysctl -w net.ipv4.ip_forward=1
    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns2 sysctl -w net.ipv4.ip_forward=1

Check Linux Routes2
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Routes Gateway    node=agent_vpp_1    namespace=ns1    ip=${IP_4}    next_hop=${IP_2}
    Wait Until Keyword Succeeds   ${WAIT_TIMEOUT}   ${SYNC_SLEEP}    Check Linux Routes Gateway    node=agent_vpp_1    namespace=ns3    ip=${IP_1}    next_hop=${IP_3}

    Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_4}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_1}


Remove VPP Nodes
    Remove All Nodes
    Sleep   ${RESYNC_SLEEP}

Start VPP1 Again
    Add Agent VPP Node    agent_vpp_1
    Sleep    ${RESYNC_WAIT}

Check Linux Interfaces On VPP1 After Resync
    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns1 ip a
    Should Contain    ${out}    ns1_veth1_linux

    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns2 ip a
    Should Contain    ${out}    ns2_veth2_linux
    Should Contain    ${out}    ns2_veth3_linux

    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns3 ip a
    Should Contain    ${out}    ns3_veth3_linux

    Check Linux Interfaces    node=agent_vpp_1    namespace=ns1    interface=ns1_veth1
    Check Linux Interfaces    node=agent_vpp_1    namespace=ns2    interface=ns2_veth2

    Check Linux Interfaces    node=agent_vpp_1    namespace=ns2    interface=ns2_veth3
    Check Linux Interfaces    node=agent_vpp_1    namespace=ns3    interface=ns3_veth3

    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns1 ip route add ${IP_4}/32 via ${IP_2}
    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns3 ip route add ${IP_1}/32 via ${IP_3}

    # IP FORWARDING
    # Advice from Andrej Marcinek
    # 9, enable IP forwarding
    # sudo sysctl -w net.ipv4.ip_forward=1
    # 10, do 'sudo nano /etc/sysctl.conf' and uncomment line:
    # net.ipv4.ip_forward=1
    #https://unix.stackexchange.com/questions/292801/routing-between-linux-namespaces
    #ip netns exec $NS_MID sysctl -w net.ipv4.ip_forward=1
    ${out}=    Execute In Container    agent_vpp_1    ip netns exec ns2 sysctl -w net.ipv4.ip_forward=1


Try to Ping6 among namespaces
    Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_2}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns2    ip=${IP_1}

    Ping6 in namespace    node=agent_vpp_1    namespace=ns2    ip=${IP_4}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_3}

    Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_3}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_2}

    Ping6 in namespace    node=agent_vpp_1    namespace=ns1    ip=${IP_4}
    Ping6 in namespace    node=agent_vpp_1    namespace=ns3    ip=${IP_1}

*** Keywords ***
Check Linux Interfaces
    [Arguments]    ${node}    ${namespace}    ${interface}
    ${out}=    Execute In Container    ${node}    ip netns exec ${namespace} ip a
    Should Contain    ${out}    ${interface}

Check Linux Routes
    [Arguments]    ${node}    ${namespace}    ${ip}
    ${out}=    Execute In Container    ${node}    ip netns exec ${namespace} ip route show
    Should Contain    ${out}    ${ip} via

Check Linux Routes Gateway
    [Arguments]    ${node}    ${namespace}    ${ip}    ${next_hop}=${EMPTY}
    ${out}=    Execute In Container    ${node}    ip netns exec ${namespace} ip route show
    Should Contain    ${out}    ${ip} via ${next_hop}

Check Linux Default Routes
    [Arguments]    ${node}    ${namespace}    ${next_hop}
    ${out}=    Execute In Container    ${node}    ip netns exec ${namespace} ip route show
    Should Contain    ${out}    default via ${next_hop}

Check Linux Routes Metric
    [Arguments]    ${node}    ${namespace}    ${ip}    ${metric}
    ${out}=    Execute In Container    ${node}    ip netns exec ${namespace} ip route show
    Should Match Regexp    ${out}    ${ip} via.*metric ${metric}\\s

Check Removed Linux Route
    [Arguments]    ${node}    ${namespace}    ${ip}
    ${out}=    Execute In Container    ${node}    ip netns exec ${namespace} ip route show
    Should Not Contain    ${out}    ${ip} via

Ping6 in namespace
    [Arguments]    ${node}    ${namespace}    ${ip}
    ${out}=    Execute In Container    ${node}    ip netns exec ${namespace} ping6 -c 5 ${ip}
    Should Contain     ${out}    from ${ip}
    Should Not Contain    ${out}    100% packet loss

TestSetup
    Make Datastore Snapshots    ${TEST_NAME}_test_setup

TestTeardown
    Make Datastore Snapshots    ${TEST_NAME}_test_teardown