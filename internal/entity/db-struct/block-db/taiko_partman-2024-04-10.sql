
SELECT partman.create_parent(p_parent_table => 'public.block_data_decode',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'daily',
                             p_premake => 30,
                             p_start_partition => '2023-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.block_data_decode';



SELECT partman.create_parent(p_parent_table => 'public.aa_block_info',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'daily',
                             p_premake => 30,
                             p_start_partition => '2023-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_block_info';





SELECT partman.create_parent(p_parent_table => 'public.transaction_decode',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'daily',
                             p_premake => 30,
                             p_start_partition => '2023-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.transaction_decode';





SELECT partman.create_parent(p_parent_table => 'public.transaction_receipt_decode',
                             p_control => 'time',
                             p_type => 'native',
                             p_interval=> 'daily',
                             p_premake => 30,
                             p_start_partition => '2023-01-01');
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.transaction_receipt_decode';







SELECT partman.create_parent(
               p_parent_table=>'public.aa_transaction_info',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'daily',
               p_premake=>30,
               p_start_partition=>'2023-01-01'
       );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_transaction_info';




SELECT partman.create_parent(
               p_parent_table=>'public.aa_user_ops_calldata',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'daily',
               p_premake=>30,
               p_start_partition=>'2023-01-01'
       );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_user_ops_calldata';



SELECT partman.create_parent(
               p_parent_table=>'public.aa_user_ops_info',
               p_control=>'time',
               p_type=>'native',
               p_interval=>'daily',
               p_premake=>30,
               p_start_partition=>'2023-01-01'
       );
UPDATE partman.part_config
SET infinite_time_partitions = TRUE
WHERE parent_table = 'public.aa_user_ops_info';